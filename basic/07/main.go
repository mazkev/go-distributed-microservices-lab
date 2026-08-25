package main

import (
	"errors"
	"fmt"
)

// 1. Definisi Custom Error sebagai variabel konstan
var (
	ErrNominalTidakValid = errors.New("nominal transaksi harus lebih besar dari 0")
	ErrSaldoKurang       = errors.New("saldo pembayaran tidak mencukupi")
)

// 2. Fungsi dengan Multiple Returns: (kembalian int, error)
func ProsesPembayaran(totalTagihan int, uangBayar int) (int, error) {
	// Validasi 1: Cek apakah nominal masuk akal
	if totalTagihan <= 0 || uangBayar <= 0 {
		return 0, ErrNominalTidakValid
	}

	// Validasi 2: Cek apakah uang cukup
	if uangBayar < totalTagihan {
		kurangnya := totalTagihan - uangBayar
		// fmt.Errorf digunakan untuk menyisipkan pesan dinamis
		return 0, fmt.Errorf("%w (kurang Rp%d)", ErrSaldoKurang, kurangnya)
	}

	// Skenario Sukses: kembalikan sisa uang dan error bernilai nil (tidak ada error)
	kembalian := uangBayar - totalTagihan
	return kembalian, nil
}

func main() {
	fmt.Println("=== MATERI 07: FUNCTIONS & IDIOMATIC ERROR HANDLING ===")

	tagihan := 45000

	// 🧪 Uji Coba Kasus A: Uang Kurang
	fmt.Println("\n--- Skenario 1: Pembayaran Kurang ---")
	kembalian, err := ProsesPembayaran(tagihan, 30000)
	if err != nil {
		fmt.Printf("❌ Transaksi Gagal: %v\n", err)
	} else {
		fmt.Printf("✅ Transaksi Berhasil! Kembalian: Rp%d\n", kembalian)
	}

	// 🧪 Uji Coba Kasus B: Nominal Negatif / 0
	fmt.Println("\n--- Skenario 2: Nominal Tidak Valid ---")
	kembalian, err = ProsesPembayaran(tagihan, -10000)
	if err != nil {
		fmt.Printf("❌ Transaksi Gagal: %v\n", err)
	} else {
		fmt.Printf("✅ Transaksi Berhasil! Kembalian: Rp%d\n", kembalian)
	}

	// 🧪 Uji Coba Kasus C: Pembayaran Sukses
	fmt.Println("\n--- Skenario 3: Pembayaran Pas / Lebih ---")
	kembalian, err = ProsesPembayaran(tagihan, 50000)
	if err != nil {
		fmt.Printf("❌ Transaksi Gagal: %v\n", err)
	} else {
		fmt.Printf("✅ Transaksi Berhasil! Kembalian: Rp%d\n", kembalian)
	}
}
