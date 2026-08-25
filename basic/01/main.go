package main

import "fmt"

func main() {
	fmt.Println("=== MATERI 01: VARIABEL, OPERATOR & KONDISI IF-ELSE ===")

	// 1. Deklarasi Variabel
	namaMenu := "Kopi Susu Gula Aren"
	harga := 18000
	jumlahBeli := 2
	uangBayar := 50000

	// 2. Perhitungan Matematika Dasar
	totalBayar := harga * jumlahBeli

	// 3. Menampilkan Ringkasan Awal
	fmt.Println("Menu      :", namaMenu)
	fmt.Println("Harga     : Rp", harga)
	fmt.Println("Jumlah    :", jumlahBeli)
	fmt.Println("Subtotal  : Rp", totalBayar)

	// 4. Logika Percabangan Diskon (if)
	if totalBayar >= 30000 {
		totalBayar = totalBayar - 5000
		fmt.Println("Diskon    : Rp 5000 (Promo > Rp30.000)")
	}

	fmt.Println("----------------------------------------")
	fmt.Println("Total Akhir : Rp", totalBayar)
	fmt.Println("Uang Bayar  : Rp", uangBayar)

	// 5. Logika Pengecekan Kembalian (if - else)
	if uangBayar >= totalBayar {
		kembalian := uangBayar - totalBayar
		fmt.Println("Kembalian   : Rp", kembalian)
		fmt.Println("Status      : ✅ PEMBAYARAN LUNAS")
	} else {
		kekurangan := totalBayar - uangBayar
		fmt.Println("Kekurangan  : Rp", kekurangan)
		fmt.Println("Status      : ❌ UANG TIDAK CUKUP")
	}
}
