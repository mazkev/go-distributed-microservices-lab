package main

import "fmt"

type Mahasiswa struct {
	Nama  string
	Nilai float64
}

// 1. Fungsi Pass by VALUE (Tanpa Pointer):
// Go membuat SALINAN baru. Perubahan di sini TIDAK mempengaruhi data asli!
func UpdateNilaiSalah(m Mahasiswa, nilaiBaru float64) {
	m.Nilai = nilaiBaru
	fmt.Println("  [Di dalam UpdateNilaiSalah] Nilai copy berubah jadi:", m.Nilai)
}

// 2. Fungsi Pass by POINTER (*Mahasiswa):
// Menerima ALAMAT MEMORI asli. Perubahan langsung merubah data aslinya!
func UpdateNilaiBenar(m *Mahasiswa, nilaiBaru float64) {
	m.Nilai = nilaiBaru
	fmt.Println("  [Di dalam UpdateNilaiBenar] Nilai asli berhasil diubah jadi:", m.Nilai)
}

// 3. Pointer Receiver pada Method Struct:
func (m *Mahasiswa) TambahBonus(bonus float64) {
	m.Nilai += bonus
}

func main() {
	fmt.Println("=== MATERI 06: POINTERS & MEMORY MUTATION ===")

	// A. Eksperimen Nilai Angka & Pointer Dasar
	x := 10
	var ptrX *int = &x // ptrX menyimpan alamat memori dari variabel x

	fmt.Printf("Nilai x           : %d\n", x)
	fmt.Printf("Alamat memori x   : %p\n", &x)
	fmt.Printf("Isi pointer ptrX  : %p\n", ptrX)
	fmt.Printf("Dereference (*ptr): %d\n", *ptrX)

	// Mengubah nilai x lewat pointer
	*ptrX = 25
	fmt.Printf("Nilai x setelah *ptrX = 25 : %d\n", x)
	fmt.Println("--------------------------------------------------")

	// B. Eksperimen Struct dengan Pass-by-Value vs Pass-by-Pointer
	mhs := Mahasiswa{Nama: "Kevin", Nilai: 70.0}
	fmt.Println("Nilai Awal Mahasiswa :", mhs.Nilai)

	// Uji Coba 1: Pakai Pass by Value (Gagal Ubah Data Asli)
	fmt.Println("\n1. Memanggil UpdateNilaiSalah(mhs, 90.0)...")
	UpdateNilaiSalah(mhs, 90.0)
	fmt.Println("   Hasil Nilai Mahasiswa Asli:", mhs.Nilai, "(TETAP 70!)")

	// Uji Coba 2: Pakai Pass by Pointer &mhs (Berhasil Ubah Data Asli)
	fmt.Println("\n2. Memanggil UpdateNilaiBenar(&mhs, 90.0)...")
	UpdateNilaiBenar(&mhs, 90.0)
	fmt.Println("   Hasil Nilai Mahasiswa Asli:", mhs.Nilai, "(BERUBAH JADI 90!)")

	// Uji Coba 3: Memanggil Pointer Receiver Method
	fmt.Println("\n3. Memanggil method mhs.TambahBonus(5.0)...")
	mhs.TambahBonus(5.0)
	fmt.Println("   Hasil Nilai Mahasiswa Akhir:", mhs.Nilai, "(MENJADI 95!)")
}
