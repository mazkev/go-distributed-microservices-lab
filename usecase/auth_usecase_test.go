package usecase_test

import (
	"errors"
	"testing"

	"gotest/domain"
	"gotest/usecase"
	"gotest/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// 1. MOCK USER REPOSITORY
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(user *domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) FindByEmail(email string) (*domain.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) FindByID(id uint) (*domain.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

// 2. UNIT TESTS

func TestAuthUsecase_Register(t *testing.T) {
	t.Run("success register new user", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		// Email belum pernah ada
		mockRepo.On("FindByEmail", "kevin@test.com").Return(nil, errors.New("not found"))
		mockRepo.On("Create", mock.AnythingOfType("*domain.User")).Return(nil)

		uc := usecase.NewAuthUsecase(mockRepo, nil)

		req := domain.RegisterRequest{
			Email:    "kevin@test.com",
			Password: "password123",
			FullName: "Kevin Pratama",
		}

		user, err := uc.Register(req)

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, "kevin@test.com", user.Email)
		assert.Equal(t, "Kevin Pratama", user.FullName)
		mockRepo.AssertExpectations(t)
	})

	t.Run("failed when email already registered", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		existingUser := &domain.User{ID: 1, Email: "kevin@test.com"}
		mockRepo.On("FindByEmail", "kevin@test.com").Return(existingUser, nil)

		uc := usecase.NewAuthUsecase(mockRepo, nil)

		req := domain.RegisterRequest{
			Email:    "kevin@test.com",
			Password: "password123",
			FullName: "Kevin Pratama",
		}

		user, err := uc.Register(req)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, "email sudah terdaftar", err.Error())
		mockRepo.AssertExpectations(t)
	})
}

func TestAuthUsecase_Login(t *testing.T) {
	hashedPassword, _ := utils.HashPassword("password123")

	t.Run("success login with correct credentials", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		user := &domain.User{
			ID:       1,
			Email:    "kevin@test.com",
			Password: hashedPassword,
			FullName: "Kevin Pratama",
		}
		mockRepo.On("FindByEmail", "kevin@test.com").Return(user, nil)

		uc := usecase.NewAuthUsecase(mockRepo, nil)

		res, err := uc.Login(domain.LoginRequest{
			Email:    "kevin@test.com",
			Password: "password123",
		})

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.NotEmpty(t, res.Token)
		assert.Equal(t, "kevin@test.com", res.User.Email)
		mockRepo.AssertExpectations(t)
	})

	t.Run("failed login with wrong password", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		user := &domain.User{
			ID:       1,
			Email:    "kevin@test.com",
			Password: hashedPassword,
		}
		mockRepo.On("FindByEmail", "kevin@test.com").Return(user, nil)

		uc := usecase.NewAuthUsecase(mockRepo, nil)

		res, err := uc.Login(domain.LoginRequest{
			Email:    "kevin@test.com",
			Password: "wrongpassword",
		})

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, "email atau password salah", err.Error())
		mockRepo.AssertExpectations(t)
	})
}
