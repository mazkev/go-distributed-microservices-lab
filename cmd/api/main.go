package main

import (
	"log"
	"log/slog"

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
	// 1. Inisialisasi JSON Structured Logger (log/slog)
	utils.InitLogger()
	slog.Info("Menginisialisasi layanan Backend API...")

	// 2. Inisialisasi Database (SQLite file: app.db)
	db, err := gorm.Open(sqlite.Open("app.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Gagal koneksi database: %v", err)
	}

	// Auto-Migrate tabel database
	err = db.AutoMigrate(
		&domain.User{},
		&domain.Product{},
		&domain.Wallet{},
		&domain.WalletTransaction{},
	)
	if err != nil {
		log.Fatalf("Gagal migrasi database: %v", err)
	}
	slog.Info("Database SQLite berhasil dimigrasi (Users, Products, Wallets, Transactions)")

	// 3. Inisialisasi Redis Cache (Graceful Fallback)
	utils.InitRedis("localhost:6379", "", 0)

	// 4. Inisialisasi Asynchronous Event Broker & Background Consumers
	broker := events.NewAsyncChannelBroker(100, 3)
	broker.Subscribe("user.registered", events.EmailNotificationConsumer)
	broker.Subscribe("user.registered", events.AuditLogConsumer)
	broker.Subscribe("product.created", events.AuditLogConsumer)
	defer broker.Close()
	slog.Info("Asynchronous Event Broker & Worker Pool aktif")

	// 5. Dependency Injection
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

	// 6. Router Setup dengan Tracing Middleware
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(http.RequestIDMiddleware())        // Menyisipkan unique X-Request-ID
	r.Use(http.StructuredLoggerMiddleware()) // Mencatat log format JSON

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

	slog.Info("🚀 Clean Architecture API berjalan", slog.String("port", ":8080"), slog.String("env", "production-ready"))
	r.Run(":8080")
}
