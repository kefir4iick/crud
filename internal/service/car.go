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

func (s *carService) GetByID(ctx context.Context, id string) (*domain.Car, error) {
	if id == "" {
		return nil, errors.New("id is required")
	}

	car, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get car: %w", err)
	}

	return car, nil
}

func (s *carService) GetAll(ctx context.Context, limit, offset int) ([]domain.Car, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	cars, err := s.repo.GetAll(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get cars: %w", err)
	}

	return cars, nil
}

func (s *carService) Update(ctx context.Context, id string, input domain.UpdateCarInput) (*domain.Car, error) {
	if id == "" {
		return nil, errors.New("id is required")
	}

	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("car not found: %w", err)
	}
	if existing == nil {
		return nil, domain.ErrCarNotFound
	}

	if input.Make != nil {
		if *input.Make == "" {
			return nil, errors.New("make cannot be empty")
		}
		if len(*input.Make) > 255 {
			return nil, errors.New("make must be less than 255 characters")
		}
		existing.Make = *input.Make
	}

	if input.Model != nil {
		if *input.Model == "" {
			return nil, errors.New("model cannot be empty")
		}
		existing.Model = *input.Model
	}

	if input.Year != nil {
		if *input.Year < 1900 {
			return nil, errors.New("year must be >= 1900")
		}
		existing.Year = *input.Year
	}

	if input.Price != nil {
		if *input.Price <= 0 {
			return nil, errors.New("price must be positive")
		}
		existing.Price = *input.Price
	}

	updated, err := s.repo.Update(ctx, id, *existing)
	if err != nil {
		return nil, fmt.Errorf("failed to update car: %w", err)
	}

	return updated, nil
}
