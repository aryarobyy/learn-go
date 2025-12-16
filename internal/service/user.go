package service

import (
	"auth/internal/model"
	"auth/internal/repository"
	"context"
	"fmt"
)

type UserService interface {
	GetById(ctx context.Context, id int) (*model.User, error)
	GetByUsername(ctx context.Context, username string) (*model.User, error)
}

type userService struct {
	repo repository.UserRepo
}

func NewUserService(repo repository.UserRepo) UserService {
	return &userService{repo: repo}
}

func (h *userService) GetById(ctx context.Context, id int) (*model.User, error) {
	res, err := h.repo.GetById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed getting user: %w", err)
	}
	return res, nil
}

func (h *userService) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	res, err := h.repo.GetByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("failed getting user: %w", err)
	}
	return res, nil
}
