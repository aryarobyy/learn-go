package service

import (
	"auth/internal/helper"
	"auth/internal/model"
	"auth/internal/repository"
	"context"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Create(ctx context.Context, user model.User) (*model.User, error)
	Login(ctx context.Context, user model.User) (*model.User, string, error)
}

type authService struct {
	repo     repository.AuthRepo
	userRepo repository.UserRepo
}

func NewAuthService(repo repository.AuthRepo, userRepo repository.UserRepo) AuthService {
	return &authService{repo: repo, userRepo: userRepo}
}

func (h *authService) Create(ctx context.Context, user model.User) (*model.User, error) {
	if user.Username == "" && user.Password == "" {
		return nil, fmt.Errorf("Username or password cannot be nul")
	}

	if !helper.IsValidName(user.Name) {
		return nil, fmt.Errorf("Name cannot contain name")
	}

	password := user.Password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hashing password", err)
	}

	registerData := model.User{
		Name:     user.Name,
		Password: string(hashedPassword),
		Username: user.Username,
	}

	res, err := h.repo.Create(ctx, registerData)
	if err != nil {
		return nil, fmt.Errorf("failed create user: %w", err)
	}
	return res, nil
}

func (h *authService) Login(ctx context.Context, user model.User) (*model.User, string, error) {
	res, err := h.userRepo.GetByUsername(ctx, user.Username)
	if err != nil {
		return nil, "", fmt.Errorf("user not found", err)
	}

	passErr := bcrypt.CompareHashAndPassword([]byte(res.Password), []byte(user.Password))
	if passErr != nil {
		return nil, "", fmt.Errorf("wrong password")
	}

	token, err := helper.CreateSessionToken(*res)
	if err != nil {
		return nil, "", fmt.Errorf("user not found", err)
	}
	return res, token, nil
}
