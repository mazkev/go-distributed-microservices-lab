package usecase

import (
	"fmt"
	"time"

	"gotest/domain"
	"gotest/repository"

	"gorm.io/gorm"
)

type WalletUsecase interface {
	GetWallet(userID uint) (*domain.Wallet, error)
	TopUp(userID uint, amount float64) (*domain.Wallet, error)
	Transfer(fromUserID uint, toUserID uint, amount float64) error
}

type walletUsecase struct {
	walletRepo repository.WalletRepository
}

func NewWalletUsecase(walletRepo repository.WalletRepository) WalletUsecase {
	return &walletUsecase{walletRepo: walletRepo}
}

func (u *walletUsecase) GetWallet(userID uint) (*domain.Wallet, error) {
	wallet, err := u.walletRepo.FindByUserID(userID)
	if err != nil {
		// Jika belum punya wallet, buatkan otomatis dengan saldo 0
		newWallet := &domain.Wallet{
			UserID:  userID,
			Balance: 0,
		}
		if err := u.walletRepo.Create(newWallet); err != nil {
			return nil, err
		}
		return newWallet, nil
	}
	return wallet, nil
}

func (u *walletUsecase) TopUp(userID uint, amount float64) (*domain.Wallet, error) {
	if amount <= 0 {
		return nil, domain.ErrInvalidTransferAmt
	}

	wallet, err := u.GetWallet(userID)
	if err != nil {
		return nil, err
	}

	wallet.Balance += amount
	if err := u.walletRepo.UpdateBalance(u.walletRepo.GetDB(), wallet); err != nil {
		return nil, err
	}

	return wallet, nil
}

// Transfer mengeksekusi transfer saldo secara ATOMIK dengan ACID Transaction & Row Lock
func (u *walletUsecase) Transfer(fromUserID uint, toUserID uint, amount float64) error {
	if fromUserID == toUserID {
		return domain.ErrSameAccountTransfer
	}
	if amount <= 0 {
		return domain.ErrInvalidTransferAmt
	}

	db := u.walletRepo.GetDB()

	// 🔒 MULAI DATABASE TRANSACTION (ACID)
	return db.Transaction(func(tx *gorm.DB) error {
		// 1. Kunci Baris Pengirim (SELECT FOR UPDATE)
		senderWallet, err := u.walletRepo.FindByUserIDWithLock(tx, fromUserID)
		if err != nil {
			return fmt.Errorf("pengirim: %w", domain.ErrWalletNotFound)
		}

		// 2. Validasi Saldo Pengirim
		if senderWallet.Balance < amount {
			return domain.ErrInsufficientBalance
		}

		// 3. Kunci Baris Penerima (SELECT FOR UPDATE)
		receiverWallet, err := u.walletRepo.FindByUserIDWithLock(tx, toUserID)
		if err != nil {
			return fmt.Errorf("penerima: %w", domain.ErrWalletNotFound)
		}

		// 4. Mutasi Saldo
		senderWallet.Balance -= amount
		receiverWallet.Balance += amount

		// 5. Update Database dalam Transaksi
		if err := u.walletRepo.UpdateBalance(tx, senderWallet); err != nil {
			return err // Otomatis Rollback jika error
		}

		if err := u.walletRepo.UpdateBalance(tx, receiverWallet); err != nil {
			return err // Otomatis Rollback jika error
		}

		// 6. Catat Histori Transaksi
		txRecord := &domain.WalletTransaction{
			FromWalletID: senderWallet.ID,
			ToWalletID:   receiverWallet.ID,
			Amount:       amount,
			Status:       "SUCCESS",
			CreatedAt:    time.Now(),
		}
		if err := u.walletRepo.CreateTransactionRecord(tx, txRecord); err != nil {
			return err // Otomatis Rollback jika error
		}

		// Jika return nil, GORM akan otomatis mengeksekusi COMMIT
		return nil
	})
}
