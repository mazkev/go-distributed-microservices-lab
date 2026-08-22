package main

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// 1. DATA MODEL & STRUCT TAGS
type Product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name" binding:"required,min=3"` // Validasi: wajib & min 3 karakter
	Price float64 `json:"price" binding:"required,gt=0"` // Validasi: wajib & > 0
	Stock int     `json:"stock" binding:"gte=0"`         // Validasi: >= 0
}

// In-Memory Database (Thread-safe dengan Mutex)
type ProductStore struct {
	mu       sync.RWMutex
	products map[int]Product
	nextID   int
}

var store = ProductStore{
	products: map[int]Product{
		1: {ID: 1, Name: "MacBook Pro M3", Price: 25000000, Stock: 10},
		2: {ID: 2, Name: "Mechanical Keyboard", Price: 1500000, Stock: 25},
	},
	nextID: 3,
}

// 2. STANDARD JSON RESPONSE HELPER
type APIResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
}

// 3. CUSTOM MIDDLEWARE: Request Logger & Timer
func RequestLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		c.Next() // Lanjutkan ke handler utama

		latency := time.Since(startTime)
		statusCode := c.Writer.Status()
		fmt.Printf("[API-LOG] %s | %d | %v | %s %s\n",
			time.Now().Format("2006-01-02 15:04:05"),
			statusCode,
			latency,
			c.Request.Method,
			c.Request.URL.Path,
		)
	}
}

// 4. HANDLERS (CONTROLLER)

// GET /api/v1/products (List all products)
func getProducts(c *gin.Context) {
	store.mu.RLock()
	defer store.mu.RUnlock()

	var list []Product
	for _, p := range store.products {
		list = append(list, p)
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    list,
	})
}

// GET /api/v1/products/:id (Get single product)
func getProductByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{Success: false, Error: "ID harus berupa angka"})
		return
	}

	store.mu.RLock()
	product, exists := store.products[id]
	store.mu.RUnlock()

	if !exists {
		c.JSON(http.StatusNotFound, APIResponse{Success: false, Error: "Produk tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, APIResponse{Success: true, Data: product})
}

// POST /api/v1/products (Create product with validation)
func createProduct(c *gin.Context) {
	var input Product
	// ShouldBindJSON otomatis memvalidasi tag `binding:"required,min=3"`
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   fmt.Sprintf("Validasi gagal: %s", err.Error()),
		})
		return
	}

	store.mu.Lock()
	input.ID = store.nextID
	store.nextID++
	store.products[input.ID] = input
	store.mu.Unlock()

	c.JSON(http.StatusCreated, APIResponse{
		Success: true,
		Message: "Produk berhasil ditambahkan",
		Data:    input,
	})
}

// DELETE /api/v1/products/:id
func deleteProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{Success: false, Error: "ID tidak valid"})
		return
	}

	store.mu.Lock()
	defer store.mu.Unlock()

	if _, exists := store.products[id]; !exists {
		c.JSON(http.StatusNotFound, APIResponse{Success: false, Error: "Produk tidak ditemukan"})
		return
	}

	delete(store.products, id)
	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: fmt.Sprintf("Produk dengan ID %d berhasil dihapus", id),
	})
}

func main() {
	// Mode rilis: gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	// Pasang custom middleware
	router.Use(RequestLoggerMiddleware())

	// Route Grouping (Standard REST API versioning)
	v1 := router.Group("/api/v1")
	{
		v1.GET("/products", getProducts)
		v1.GET("/products/:id", getProductByID)
		v1.POST("/products", createProduct)
		v1.DELETE("/products/:id", deleteProduct)
	}

	fmt.Println(" Server REST API berjalan di http://localhost:8080")
	router.Run(":8080")
}
