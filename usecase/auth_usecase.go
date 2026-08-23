package usecase

import (
	"context"
	"errors"
	"time"

	"gotest/domain"
	"gotest/events"
	"gotest/repository"
	"gotest/utils"
)

type AuthUsecase interface {
	Register(req domain.RegisterRequest) (*domain.User, error)
	Login(req domain.LoginRequest) (*domain.LoginResponse, error)
}

type authUsecase struct {
	userRepo repository.UserRepository
	broker   events.EventBroker
}

func NewAuthUsecase(userRepo repository.UserRepository, broker events.EventBroker) AuthUsecase {
	return &authUsecase{
		userRepo: userRepo,
		broker:   broker,
	}
}

func (u *authUsecase) Register(req domain.RegisterRequest) (*domain.User, error) {
	// Cek apakah email sudah terdaftar
	existing, _ := u.userRepo.FindByEmail(req.Email)
	if existing != nil {
		return nil, errors.New("email sudah terdaftar")
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := domain.User{
		Email:    req.Email,
		Password: hashedPassword,
		FullName: req.FullName,
	}

	err = u.userRepo.Create(&user)
	if err != nil {
		return nil, err
	}

	// 📨 Publish Asynchronous Event (Non-blocking)
	if u.broker != nil {
		_ = u.broker.Publish(context.Background(), events.UserRegisteredEvent{
			UserID:    user.ID,
			Email:     user.Email,
			FullName:  user.FullName,
			Timestamp: time.Now(),
		})
	}

	return &user, nil
}

func (u *authUsecase) Login(req domain.LoginRequest) (*domain.LoginResponse, error) {
	user, err := u.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, errors.New("email atau password salah")
	}

	if !utils.CheckPassword(req.Password, user.Password) {
		return nil, errors.New("email atau password salah")
	}

	token, err := utils.GenerateToken(user.ID, user.Email)
	if err != nil {
		return nil, err
	}

	return &domain.LoginResponse{
		Token: token,
		User:  *user,
	}, nil
}
