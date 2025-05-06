package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/kefir4iick/go/internal/domain"
	"github.com/kefir4iick/go/internal/repository"
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
