package repository

import (
	"context"
	"github.com/kefir4iick/crud/internal/domain"
)

type CarRepository interface {
	Create(ctx context.Context, car domain.Car) (*domain.Car, error)
	GetByID(ctx context.Context, id string) (*domain.Car, error)
	GetAll(ctx context.Context, limit, offset int) ([]domain.Car, error)
	Update(ctx context.Context, id string, car domain.Car) (*domain.Car, error)
	Delete(ctx context.Context, id string) error
}
