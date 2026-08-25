package main

import "fmt"

func main() {
	angka := []int{12, 45, 7, 89, 24, 56, 11, 80}

	// 1. Deklarasi 2 slice kosong untuk menampung hasil pemisahan
	var angkaGenap []int
	var angkaGanjil []int

	// 2. Looping dan filter menggunakan Modulo (%) dan append()
	for _, nilai := range angka {
		if nilai%2 == 0 {
			// Jika habis dibagi 2 -> Masukkan ke slice angkaGenap
			angkaGenap = append(angkaGenap, nilai)
		} else {
			// Jika ada sisa pembagian -> Masukkan ke slice angkaGanjil
			angkaGanjil = append(angkaGanjil, nilai)
		}
	}

	// 3. Cetak Hasil Pemisahan
	fmt.Println("Daftar Angka Asli :", angka)
	fmt.Println("Angka Genap       :", angkaGenap)
	fmt.Println("Angka Ganjil      :", angkaGanjil)
	fmt.Println("Jumlah Genap      :", len(angkaGenap))
	fmt.Println("Jumlah Ganjil     :", len(angkaGanjil))
}
