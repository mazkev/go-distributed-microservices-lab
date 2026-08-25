package main

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Println("=== MATERI 03: STRING MANIPULATION & PALINDROME ===")

	kata := "Katak"
	// 1. Ubah ke huruf kecil semua agar pengecekan adil (tidak beda antara 'K' dan 'k')
	kataBersih := strings.ToLower(kata)

	var kataTerbalik string
	jumlahVokal := 0

	// 2. Looping Mundur untuk Membalik Kata
	for i := len(kataBersih) - 1; i >= 0; i-- {
		huruf := string(kataBersih[i])
		kataTerbalik += huruf

		// 3. Menghitung Jumlah Huruf Vokal
		if huruf == "a" || huruf == "i" || huruf == "u" || huruf == "e" || huruf == "o" {
			jumlahVokal++
		}
	}

	// 4. Tampilkan Hasil Analisis String
	fmt.Println("Kata Asli       :", kata)
	fmt.Println("Kata Terbalik   :", kataTerbalik)
	fmt.Println("Jumlah Karakter :", len(kata))
	fmt.Println("Jumlah Vokal    :", jumlahVokal)
	fmt.Println("----------------------------------------")

	// 5. Cek Palindrome
	if kataBersih == kataTerbalik {
		fmt.Println("Status          : ✅ Kata ini adalah PALINDROME!")
	} else {
		fmt.Println("Status          : ❌ BUKAN kata Palindrome.")
	}
}
