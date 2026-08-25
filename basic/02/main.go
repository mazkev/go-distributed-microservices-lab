package main

import "fmt"

func main() {
	fmt.Println("=== MATERI 02: SLICES, LOOPS & AGGREGATE LOGIC ===")

	// 1. Data Slice
	angka := []int{12, 45, 7, 89, 24, 56, 11, 80}

	// 2. Variabel Penampung
	total := 0
	max := angka[0]
	min := angka[0]
	var angkaGenap []int
	var angkaGanjil []int

	// 3. Looping Tunggal untuk Menyelesaikan Semua Operasi
	for _, nilai := range angka {
		// A. Hitung Total
		total += nilai

		// B. Cari Nilai Maksimum
		if nilai > max {
			max = nilai
		}

		// C. Cari Nilai Minimum
		if nilai < min {
			min = nilai
		}

		// D. Filter Genap vs Ganjil menggunakan Modulo (%) dan append()
		if nilai%2 == 0 {
			angkaGenap = append(angkaGenap, nilai)
		} else {
			angkaGanjil = append(angkaGanjil, nilai)
		}
	}

	// 4. Hitung Rata-rata dengan konversi float64
	rataRata := float64(total) / float64(len(angka))

	// 5. Cetak Hasil Analisis Data
	fmt.Println("Daftar Angka  :", angka)
	fmt.Println("Total Nilai   :", total)
	fmt.Println("Angka Terbesar:", max)
	fmt.Println("Angka Terkecil:", min)
	fmt.Printf("Rata-rata     : %.2f\n", rataRata)
	fmt.Println("----------------------------------------")
	fmt.Println("Angka Genap   :", angkaGenap)
	fmt.Println("Angka Ganjil  :", angkaGanjil)
}
