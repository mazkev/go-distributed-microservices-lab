package main

import (
	"fmt"
	"log"

	"gotest/delivery/http"
	"gotest/domain"
	"gotest/events"
	"gotest/repository"
	"gotest/usecase"
	"gotest/utils"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func main() {
	// 1. Inisialisasi Database (SQLite file: app.db)
	db, err := gorm.Open(sqlite.Open("app.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Gagal koneksi database: %v", err)
	}

	// Auto-Migrate semua tabel database secara otomatis
	err = db.AutoMigrate(
		&domain.User{},
		&domain.Product{},
		&domain.Wallet{},
		&domain.WalletTransaction{},
	)
	if err != nil {
		log.Fatalf("Gagal migrasi database: %v", err)
	}
	fmt.Println("✅ Database SQLite berhasil dimigrasi (Users, Products, Wallets, Transactions)!")

	// 2. Inisialisasi Redis Cache (Graceful Fallback)
	utils.InitRedis("localhost:6379", "", 0)

	// 3. Inisialisasi Asynchronous Event Broker & Background Consumers
	broker := events.NewAsyncChannelBroker(100, 3)
	broker.Subscribe("user.registered", events.EmailNotificationConsumer)
	broker.Subscribe("user.registered", events.AuditLogConsumer)
	broker.Subscribe("product.created", events.AuditLogConsumer)
	defer broker.Close()
	fmt.Println("✅ Asynchronous Event Broker & Worker Pool aktif!")

	// 4. Dependency Injection
	// Repositories
	userRepo := repository.NewUserRepository(db)
	productRepo := repository.NewProductRepository(db)
	walletRepo := repository.NewWalletRepository(db)

	// Usecases
	authUsecase := usecase.NewAuthUsecase(userRepo, broker)
	productUsecase := usecase.NewProductUsecase(productRepo)
	walletUsecase := usecase.NewWalletUsecase(walletRepo)

	// Handlers
	authHandler := http.NewAuthHandler(authUsecase)
	productHandler := http.NewProductHandler(productUsecase)
	walletHandler := http.NewWalletHandler(walletUsecase)

	// 5. Router Setup
	r := gin.Default()

	api := r.Group("/api/v1")
	{
		// Public Routes (Auth)
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		// Public Routes (Baca Produk)
		api.GET("/products", productHandler.GetAll)
		api.GET("/products/:id", productHandler.GetByID)

		// Protected Routes (Wajib Login & kirim Header Authorization: Bearer <token>)
		protected := api.Group("/")
		protected.Use(http.AuthMiddleware())
		{
			// Products
			protected.POST("/products", productHandler.Create)
			protected.DELETE("/products/:id", productHandler.Delete)

			// Wallets & Transactions (Fintech Module)
			protected.GET("/wallets/me", walletHandler.GetMyWallet)
			protected.POST("/wallets/topup", walletHandler.TopUp)
			protected.POST("/wallets/transfer", walletHandler.Transfer)
		}
	}

	fmt.Println("🚀 Clean Architecture API berjalan di http://localhost:8080")
	r.Run(":8080")
}
