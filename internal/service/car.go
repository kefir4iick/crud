package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/kefir4iick/crud/internal/domain"
	"github.com/kefir4iick/crud/internal/repository"
)

type CarService interface {
	Create(ctx context.Context, input domain.Car) (*domain.Car, error)
	GetByID(ctx context.Context, id string) (*domain.Car, error)
	GetAll(ctx context.Context, limit, offset int) ([]domain.Car, error)
	Update(ctx context.Context, id string, input domain.UpdateCarInput) (*domain.Car, error)
	Delete(ctx context.Context, id string) error
}

type carService struct {
	repo repository.CarRepository
}

func NewCarService(repo repository.CarRepository) CarService {
	return &carService{repo: repo}
}

func (s *carService) Create(ctx context.Context, input domain.Car) (*domain.Car, error) {
	if input.Make == "" {
		return nil, errors.New("make is required")
	}
	if len(input.Make) > 255 {
		return nil, errors.New("make must be less than 255 characters")
	}
	if input.Model == "" {
		return nil, errors.New("model is required")
	}
	if input.Year < 1900 {
		return nil, errors.New("year must be >= 1900")
	}
	if input.Price <= 0 {
		return nil, errors.New("price must be positive")
	}

	car, err := s.repo.Create(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to create car: %w", err)
	}

	return car, nil
}
