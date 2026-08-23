package repository

import (
	"gotest/domain"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type WalletRepository interface {
	GetDB() *gorm.DB
	Create(wallet *domain.Wallet) error
	FindByUserID(userID uint) (*domain.Wallet, error)
	FindByUserIDWithLock(tx *gorm.DB, userID uint) (*domain.Wallet, error)
	UpdateBalance(tx *gorm.DB, wallet *domain.Wallet) error
	CreateTransactionRecord(tx *gorm.DB, record *domain.WalletTransaction) error
}

type walletRepository struct {
	db *gorm.DB
}

func NewWalletRepository(db *gorm.DB) WalletRepository {
	return &walletRepository{db: db}
}

func (r *walletRepository) GetDB() *gorm.DB {
	return r.db
}

func (r *walletRepository) Create(wallet *domain.Wallet) error {
	return r.db.Create(wallet).Error
}

func (r *walletRepository) FindByUserID(userID uint) (*domain.Wallet, error) {
	var wallet domain.Wallet
	err := r.db.Where("user_id = ?", userID).First(&wallet).Error
	if err != nil {
		return nil, err
	}
	return &wallet, nil
}

// FindByUserIDWithLock menggunakan SELECT FOR UPDATE (Pessimistic Row Locking)
func (r *walletRepository) FindByUserIDWithLock(tx *gorm.DB, userID uint) (*domain.Wallet, error) {
	var wallet domain.Wallet
	// Mengunci baris data wallet di DB agar thread lain menunggu hingga transaksi selesai
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("user_id = ?", userID).
		First(&wallet).Error
	if err != nil {
		return nil, err
	}
	return &wallet, nil
}

func (r *walletRepository) UpdateBalance(tx *gorm.DB, wallet *domain.Wallet) error {
	return tx.Model(wallet).Update("balance", wallet.Balance).Error
}

func (r *walletRepository) CreateTransactionRecord(tx *gorm.DB, record *domain.WalletTransaction) error {
	return tx.Create(record).Error
}
