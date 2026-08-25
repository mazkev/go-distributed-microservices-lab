package main

import "fmt"

func main() {
	fmt.Println("=== MATERI 04: MAPS & FREQUENCY COUNTER ===")

	// 1. Kalimat yang ingin dihitung frekuensi kemunculan karakternya
	teks := "golang programming"

	// 2. Deklarasi Map (Key: string huruf, Value: int jumlah kemunculan)
	frekuensi := make(map[string]int)

	// 3. Looping setiap karakter dan hitung kemunculannya
	for _, karakter := range teks {
		huruf := string(karakter)

		// Lewatkan spasi (hanya hitung huruf)
		if huruf == " " {
			continue
		}

		// Tambahkan hitungan untuk huruf ini
		frekuensi[huruf]++
	}

	// 4. Cetak Tabel Frekuensi Huruf
	fmt.Println("Teks Asli:", teks)
	fmt.Println("----------------------------------------")
	fmt.Println("Huruf | Frekuensi")
	fmt.Println("----------------------------------------")
	for huruf, jumlah := range frekuensi {
		fmt.Printf("  %s   |   %d kali\n", huruf, jumlah)
	}

	fmt.Println("----------------------------------------")
	// 5. Pengecekan Keberadaan Key di Map (Idiom `value, exists`)
	cariHuruf := "g"
	if jumlah, ada := frekuensi[cariHuruf]; ada {
		fmt.Printf("🔍 Huruf '%s' ditemukan sebanyak %d kali.\n", cariHuruf, jumlah)
	} else {
		fmt.Printf("🔍 Huruf '%s' tidak ditemukan dalam teks.\n", cariHuruf)
	}
}
