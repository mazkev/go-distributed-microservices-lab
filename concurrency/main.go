package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// 1. MUTEX: Menghindari Race Condition
type SafeCounter struct {
	mu    sync.Mutex
	count int
}

func (c *SafeCounter) Increment(wg *sync.WaitGroup) {
	defer wg.Done() // Pastikan wg.Done() dipanggil saat fungsi selesai

	c.mu.Lock()         // Kunci akses: Hanya 1 goroutine yang boleh modifikasi
	c.count++
	c.mu.Unlock()       // Buka kembali kuncinya
}

// 2. CHANNEL & WORKER POOL PATTERN
// Menerima jobs dari channel dan mengirimkan hasil ke channel result
func worker(id int, jobs <-chan int, results chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	for j := range jobs {
		fmt.Printf("👷 Worker #%d mulai memproses job #%d\n", id, j)
		time.Sleep(100 * time.Millisecond) // Simulasi beban kerja
		results <- j * 2                   // Kirim hasil pemrosesan (kali 2)
	}
}

// 3. CONTEXT: Timeout & Cancellation (Pola standar di backend API Go)
func fetchUserData(ctx context.Context) (string, error) {
	// Simulasi panggilan database atau API eksternal yang butuh waktu 300ms
	select {
	case <-time.After(300 * time.Millisecond):
		return "Data User: Kevin Pratama", nil
	case <-ctx.Done():
		// Jika context dibatalkan (karena timeout) sebelum 300ms selesai
		return "", ctx.Err()
	}
}

func main() {
	fmt.Println("=== 1. SYNC.MUTEX & WAITGROUP (Mencegah Race Condition) ===")
	counter := SafeCounter{}
	var wg sync.WaitGroup

	// Jalankan 100 goroutine secara bersamaan untuk menambah counter
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go counter.Increment(&wg)
	}

	wg.Wait() // Tunggu ke-100 goroutine selesai
	fmt.Printf("Total hitungan aman (harus 100): %d\n", counter.count)

	fmt.Println("\n=== 2. CHANNELS & WORKER POOL ===")
	numJobs := 5
	jobs := make(chan int, numJobs)
	results := make(chan int, numJobs)
	var workerWg sync.WaitGroup

	// Buat 3 Worker paralel
	for w := 1; w <= 3; w++ {
		workerWg.Add(1)
		go worker(w, jobs, results, &workerWg)
	}

	// Kirim 5 jobs ke channel
	for j := 1; j <= numJobs; j++ {
		jobs <- j
	}
	close(jobs) // Tutup channel jobs menandakan tidak ada job baru

	workerWg.Wait() // Tunggu semua worker selesai mengambil jobs
	close(results)

	// Baca hasil pemrosesan dari channel results
	for r := range results {
		fmt.Printf(" Hasil diterima: %d\n", r)
	}

	fmt.Println("\n=== 3. CONTEXT DENGAN TIMEOUT ===")
	// Kasus A: Timeout 500ms (Cukup waktu, proses 300ms berhasil)
	ctxSuccess, cancelA := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancelA()

	data, err := fetchUserData(ctxSuccess)
	if err != nil {
		fmt.Printf("❌ Gagal fetch data: %v\n", err)
	} else {
		fmt.Printf("✅ Berhasil: %s\n", data)
	}

	// Kasus B: Timeout 100ms (Terlalu cepat, proses 300ms dibatalkan otomatis)
	ctxTimeout, cancelB := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancelB()

	data, err = fetchUserData(ctxTimeout)
	if err != nil {
		fmt.Printf("⏱️ Request dibatalkan sesuai ekspektasi: %v\n", err)
	} else {
		fmt.Printf("✅ Berhasil: %s\n", data)
	}
}
