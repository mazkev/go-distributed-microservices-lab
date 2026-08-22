package main

import (
	"errors"
	"fmt"
)

// 1. DEFINISI CUSTOM ERROR
// Di Go, error adalah interface sederhana: type error interface { Error() string }
var (
	ErrInsufficientBalance = errors.New("saldo tidak mencukupi")
	ErrInvalidAmount       = errors.New("jumlah nominal tidak valid")
)

// 2. STRUCT & POINTER RECEIVER
type BankAccount struct {
	Owner   string
	Balance float64
}

// Deposit: Menggunakan pointer receiver (*BankAccount) agar perubahan saldo tersimpan di memory asli
func (b *BankAccount) Deposit(amount float64) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}
	b.Balance += amount
	return nil
}

// Withdraw: Method dengan validasi error khas Go (if err != nil)
func (b *BankAccount) Withdraw(amount float64) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}
	if amount > b.Balance {
		return ErrInsufficientBalance
	}
	b.Balance -= amount
	return nil
}

// 3. INTERFACE (Kunci arsitektur modular & testing di Go)
// Tipe data apapun yang memiliki method Pay() otomatis mengimplementasikan interface PaymentProcessor
type PaymentProcessor interface {
	Pay(amount float64) error
}

type EWallet struct {
	Provider string
	Balance  float64
}

func (e *EWallet) Pay(amount float64) error {
	if amount > e.Balance {
		return fmt.Errorf("[%s] gagal bayar: saldo kurang", e.Provider)
	}
	e.Balance -= amount
	fmt.Printf("✅ Pembayaran berhasil via %s sebesar Rp%.2f (Sisa: Rp%.2f)\n", e.Provider, amount, e.Balance)
	return nil
}

// Fungsi serbaguna yang menerima implementasi interface apapun
func ProcessCheckout(p PaymentProcessor, amount float64) {
	err := p.Pay(amount)
	if err != nil {
		fmt.Printf("❌ Transaksi gagal: %v\n", err)
		return
	}
}

func main() {
	fmt.Println("=== 1. STRUCT, METHOD & POINTER ===")
	account := BankAccount{
		Owner:   "Budi",
		Balance: 100000,
	}

	fmt.Printf("Akun: %s, Saldo Awal: Rp%.2f\n", account.Owner, account.Balance)

	// Uji Deposit
	account.Deposit(50000)
	fmt.Printf("Setelah deposit: Rp%.2f\n", account.Balance)

	// Uji Error Handling saat penarikan melebihi saldo
	err := account.Withdraw(200000)
	if err != nil {
		// Pattern standar di Go: Selalu cek if err != nil
		fmt.Printf("Penarikan gagal: %v\n", err)
	}

	fmt.Println("\n=== 2. INTERFACE & POLYMORPHISM ===")
	gopay := &EWallet{Provider: "GoPay", Balance: 75000}
	ovo := &EWallet{Provider: "OVO", Balance: 20000}

	// Kedua wallet bisa diproses oleh fungsi yang sama karena mengimplementasikan PaymentProcessor
	ProcessCheckout(gopay, 50000) // Sukses
	ProcessCheckout(ovo, 50000)   // Gagal

	fmt.Println("\n=== 3. SLICE & MAP (Koleksi Data) ===")
	// Map: Key-Value store di memory
	users := map[string]string{
		"USR01": "Andi",
		"USR02": "Siti",
	}

	// Pengecekan keberadaan key di map (Idiom `value, exists`)
	if name, exists := users["USR01"]; exists {
		fmt.Printf("Ditemukan: %s\n", name)
	}

	// Slice: Dynamic array yang sangat sering dipakai
	transactions := []float64{50000, 25000, 100000}
	for i, val := range transactions {
		fmt.Printf("Transaksi #%d: Rp%.2f\n", i+1, val)
	}
}
