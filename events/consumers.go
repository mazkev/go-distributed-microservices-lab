package events

import (
	"context"
	"fmt"
	"time"
)

// EmailNotificationConsumer mensimulasikan worker pengiriman email aktivasi/selamat datang
func EmailNotificationConsumer(ctx context.Context, event Event) error {
	e, ok := event.(UserRegisteredEvent)
	if !ok {
		return fmt.Errorf("tipe event tidak cocok untuk email consumer")
	}

	// Simulasi I/O pengiriman email (butuh waktu 150ms)
	time.Sleep(150 * time.Millisecond)

	fmt.Printf("📧 [WORKER - EMAIL] Mengirim 'Welcome Email' ke: %s (%s)\n", e.Email, e.FullName)
	return nil
}

// AuditLogConsumer mencatat semua event ke log audit keamanan
func AuditLogConsumer(ctx context.Context, event Event) error {
	switch e := event.(type) {
	case UserRegisteredEvent:
		fmt.Printf("📝 [WORKER - AUDIT] User baru terdaftar: UserID=%d, Email=%s pada %s\n",
			e.UserID, e.Email, e.Timestamp.Format("2006-01-02 15:04:05"))
	case ProductCreatedEvent:
		fmt.Printf("📝 [WORKER - AUDIT] Produk baru dibuat: ProductID=%d, Nama=%s, Harga=Rp%.2f\n",
			e.ProductID, e.Name, e.Price)
	}
	return nil
}
