package main

import "fmt"

// 1. Definisi Struct Mahasiswa
type Mahasiswa struct {
	NIM     string
	Nama    string
	Jurusan string
	Nilai   float64
}

// 2. Method Receiver: Menentukan Status Kelulusan (Nilai >= 75)
func (m Mahasiswa) ApakahLulus() bool {
	return m.Nilai >= 75.0
}

// 3. Method Receiver: Menentukan Predikat Huruf (A, B, C, D)
func (m Mahasiswa) Predikat() string {
	switch {
	case m.Nilai >= 85:
		return "A (Sangat Memuaskan)"
	case m.Nilai >= 75:
		return "B (Memuaskan)"
	case m.Nilai >= 60:
		return "C (Cukup)"
	default:
		return "D (Tidak Lulus)"
	}
}

func main() {
	fmt.Println("=== MATERI 05: STRUCTS & METHODS ===")

	// 4. Data Kumpulan Mahasiswa (Slice of Structs)
	daftarMahasiswa := []Mahasiswa{
		{NIM: "101", Nama: "Andi Pratama", Jurusan: "Informatika", Nilai: 88.5},
		{NIM: "102", Nama: "Siti Rahma", Jurusan: "Sistem Informasi", Nilai: 92.0},
		{NIM: "103", Nama: "Budi Santoso", Jurusan: "Informatika", Nilai: 68.0},
		{NIM: "104", Nama: "Dewi Lestari", Jurusan: "Teknik Komputer", Nilai: 78.5},
	}

	totalNilai := 0.0
	topMahasiswa := daftarMahasiswa[0]

	// 5. Looping & Pengolahan Data
	fmt.Println("----------------------------------------------------------------------")
	fmt.Printf("%-5s | %-15s | %-16s | %-6s | %s\n", "NIM", "Nama", "Jurusan", "Nilai", "Predikat")
	fmt.Println("----------------------------------------------------------------------")

	for _, mhs := range daftarMahasiswa {
		totalNilai += mhs.Nilai

		// Cari Nilai Tertinggi
		if mhs.Nilai > topMahasiswa.Nilai {
			topMahasiswa = mhs
		}

		// Cetak baris data dengan memanggil method
		status := "LULUS"
		if !mhs.ApakahLulus() {
			status = "REMEDIAL"
		}

		fmt.Printf("%-5s | %-15s | %-16s | %-6.1f | %s [%s]\n",
			mhs.NIM, mhs.Nama, mhs.Jurusan, mhs.Nilai, mhs.Predikat(), status)
	}

	rataRata := totalNilai / float64(len(daftarMahasiswa))

	fmt.Println("----------------------------------------------------------------------")
	fmt.Printf("Rata-rata Nilai Kelas : %.2f\n", rataRata)
	fmt.Printf("Mahasiswa Terbaik     : %s (NIM: %s) dengan Nilai %.1f\n",
		topMahasiswa.Nama, topMahasiswa.NIM, topMahasiswa.Nilai)
}
