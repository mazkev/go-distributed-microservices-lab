package events

import "time"

// Event adalah interface dasar untuk semua jenis event di sistem
type Event interface {
	EventType() string
	OccurredAt() time.Time
}

// UserRegisteredEvent dipicu saat user berhasil mendaftar
type UserRegisteredEvent struct {
	UserID    uint      `json:"user_id"`
	Email     string    `json:"email"`
	FullName  string    `json:"full_name"`
	Timestamp time.Time `json:"timestamp"`
}

func (e UserRegisteredEvent) EventType() string {
	return "user.registered"
}

func (e UserRegisteredEvent) OccurredAt() time.Time {
	return e.Timestamp
}

// ProductCreatedEvent dipicu saat produk baru ditambahkan ke katalog
type ProductCreatedEvent struct {
	ProductID uint      `json:"product_id"`
	Name      string    `json:"name"`
	Price     float64   `json:"price"`
	Timestamp time.Time `json:"timestamp"`
}

func (e ProductCreatedEvent) EventType() string {
	return "product.created"
}

func (e ProductCreatedEvent) OccurredAt() time.Time {
	return e.Timestamp
}
