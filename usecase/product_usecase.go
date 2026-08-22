package usecase

import (
	"gotest/domain"
	"gotest/repository"
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
	return &product, nil
}

func (u *productUsecase) GetAll() ([]domain.Product, error) {
	return u.productRepo.FindAll()
}

func (u *productUsecase) GetByID(id uint) (*domain.Product, error) {
	return u.productRepo.FindByID(id)
}

func (u *productUsecase) Delete(id uint) error {
	return u.productRepo.Delete(id)
}
