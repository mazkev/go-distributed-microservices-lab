package usecase

import (
	"context"
	"fmt"
	"time"

	"gotest/domain"
	"gotest/repository"
	"gotest/utils"
)

type ProductUsecase interface {
	Create(req domain.ProductRequest, userID uint) (*domain.Product, error)
	GetAll() ([]domain.Product, error)
	GetByID(id uint) (*domain.Product, error)
	Delete(id uint) error
}

type productUsecase struct {
	productRepo repository.ProductRepository
}

func NewProductUsecase(productRepo repository.ProductRepository) ProductUsecase {
	return &productUsecase{productRepo: productRepo}
}

func (u *productUsecase) Create(req domain.ProductRequest, userID uint) (*domain.Product, error) {
	product := domain.Product{
		Name:   req.Name,
		Price:  req.Price,
		Stock:  req.Stock,
		UserID: userID,
	}

	err := u.productRepo.Create(&product)
	if err != nil {
		return nil, err
	}

	// Cache Invalidation: Hapus cache list produk karena ada penambahan data baru
	if utils.Redis != nil {
		_ = utils.Redis.Del(context.Background(), "products:all")
	}

	return &product, nil
}

func (u *productUsecase) GetAll() ([]domain.Product, error) {
	cacheKey := "products:all"
	var cachedProducts []domain.Product

	// 1. Cek Cache Redis (Cache Hit)
	if utils.Redis != nil && utils.Redis.Get(context.Background(), cacheKey, &cachedProducts) {
		fmt.Println("⚡ [REDIS CACHE HIT] Mengambil list produk dari Redis")
		return cachedProducts, nil
	}

	// 2. Cache Miss: Ambil dari Database
	fmt.Println("🔍 [CACHE MISS] Query list produk dari Database...")
	products, err := u.productRepo.FindAll()
	if err != nil {
		return nil, err
	}

	// 3. Simpan hasil query ke Redis dengan TTL 2 menit
	if utils.Redis != nil {
		_ = utils.Redis.Set(context.Background(), cacheKey, products, 2*time.Minute)
	}

	return products, nil
}

func (u *productUsecase) GetByID(id uint) (*domain.Product, error) {
	cacheKey := fmt.Sprintf("product:%d", id)
	var cachedProduct domain.Product

	// 1. Cek Cache Redis (Cache Hit)
	if utils.Redis != nil && utils.Redis.Get(context.Background(), cacheKey, &cachedProduct) {
		fmt.Printf("⚡ [REDIS CACHE HIT] Mengambil produk ID %d dari Redis\n", id)
		return &cachedProduct, nil
	}

	// 2. Cache Miss: Ambil dari Database
	fmt.Printf("🔍 [CACHE MISS] Query produk ID %d dari Database...\n", id)
	product, err := u.productRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// 3. Simpan ke Redis dengan TTL 5 menit
	if utils.Redis != nil {
		_ = utils.Redis.Set(context.Background(), cacheKey, product, 5*time.Minute)
	}

	return product, nil
}

func (u *productUsecase) Delete(id uint) error {
	if err := u.productRepo.Delete(id); err != nil {
		return err
	}

	// Cache Invalidation: Hapus cache item ini & cache all
	if utils.Redis != nil {
		cacheKey := fmt.Sprintf("product:%d", id)
		_ = utils.Redis.Del(context.Background(), cacheKey, "products:all")
	}

	return nil
}
