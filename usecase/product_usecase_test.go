package usecase_test

import (
	"errors"
	"testing"

	"gotest/domain"
	"gotest/usecase"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// 1. MOCK REPOSITORY MENGGUNAKAN TESTIFY/MOCK
type MockProductRepository struct {
	mock.Mock
}

func (m *MockProductRepository) Create(product *domain.Product) error {
	args := m.Called(product)
	return args.Error(0)
}

func (m *MockProductRepository) FindAll() ([]domain.Product, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Product), args.Error(1)
}

func (m *MockProductRepository) FindByID(id uint) (*domain.Product, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Product), args.Error(1)
}

func (m *MockProductRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

// 2. UNIT TESTS (TABLE-DRIVEN TEST PATTERN)

func TestProductUsecase_Create(t *testing.T) {
	t.Run("success create product", func(t *testing.T) {
		mockRepo := new(MockProductRepository)
		// Ekspektasi: repo.Create akan dipanggil dengan argument apapun dan mengembalikan nil error
		mockRepo.On("Create", mock.AnythingOfType("*domain.Product")).Return(nil)

		uc := usecase.NewProductUsecase(mockRepo)

		req := domain.ProductRequest{
			Name:  "Laptop Gaming",
			Price: 15000000,
			Stock: 10,
		}

		result, err := uc.Create(req, 1)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "Laptop Gaming", result.Name)
		assert.Equal(t, 15000000.0, result.Price)
		assert.Equal(t, uint(1), result.UserID)

		mockRepo.AssertExpectations(t)
	})

	t.Run("failed when database error", func(t *testing.T) {
		mockRepo := new(MockProductRepository)
		mockRepo.On("Create", mock.AnythingOfType("*domain.Product")).Return(errors.New("db error"))

		uc := usecase.NewProductUsecase(mockRepo)

		req := domain.ProductRequest{
			Name:  "Laptop Gaming",
			Price: 15000000,
			Stock: 10,
		}

		result, err := uc.Create(req, 1)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "db error", err.Error())

		mockRepo.AssertExpectations(t)
	})
}

func TestProductUsecase_GetByID(t *testing.T) {
	t.Run("found product", func(t *testing.T) {
		mockRepo := new(MockProductRepository)
		expectedProduct := &domain.Product{
			ID:    1,
			Name:  "Keyboard RGB",
			Price: 500000,
		}
		mockRepo.On("FindByID", uint(1)).Return(expectedProduct, nil)

		uc := usecase.NewProductUsecase(mockRepo)
		product, err := uc.GetByID(1)

		assert.NoError(t, err)
		assert.NotNil(t, product)
		assert.Equal(t, "Keyboard RGB", product.Name)
		mockRepo.AssertExpectations(t)
	})

	t.Run("not found product", func(t *testing.T) {
		mockRepo := new(MockProductRepository)
		mockRepo.On("FindByID", uint(99)).Return(nil, errors.New("record not found"))

		uc := usecase.NewProductUsecase(mockRepo)
		product, err := uc.GetByID(99)

		assert.Error(t, err)
		assert.Nil(t, product)
		mockRepo.AssertExpectations(t)
	})
}
