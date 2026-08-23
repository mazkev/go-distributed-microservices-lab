package domain

import (
	"errors"
	"time"
)

var (
	ErrWalletNotFound      = errors.New("wallet tidak ditemukan")
	ErrInsufficientBalance = errors.New("saldo tidak mencukupi untuk transfer")
	ErrSameAccountTransfer = errors.New("tidak bisa mentransfer ke akun sendiri")
	ErrInvalidTransferAmt  = errors.New("nominal transfer harus lebih dari 0")
)

type Wallet struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"unique;not null"`
	Balance   float64   `json:"balance" gorm:"not null;default:0"`
	UpdatedAt time.Time `json:"updated_at"`
}

type WalletTransaction struct {
	ID             uint      `json:"id" gorm:"primaryKey"`
	FromWalletID   uint      `json:"from_wallet_id"`
	ToWalletID     uint      `json:"to_wallet_id"`
	Amount         float64   `json:"amount"`
	Status         string    `json:"status"` // SUCCESS, FAILED
	CreatedAt      time.Time `json:"created_at"`
}

type TransferRequest struct {
	ToUserID uint    `json:"to_user_id" binding:"required"`
	Amount   float64 `json:"amount" binding:"required,gt=0"`
}

type TopUpRequest struct {
	Amount float64 `json:"amount" binding:"required,gt=0"`
}
