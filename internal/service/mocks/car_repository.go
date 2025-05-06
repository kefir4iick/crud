package mocks

import (
	"context"
	"github.com/kefir4iick/crud/internal/domain"
	"github.com/stretchr/testify/mock"
)

type CarRepository struct {
	mock.Mock
}

func (m *CarRepository) Create(ctx context.Context, car domain.Car) (*domain.Car, error) {
	args := m.Called(ctx, car)
	return args.Get(0).(*domain.Car), args.Error(1)
}

func (m *CarRepository) GetByID(ctx context.Context, id string) (*domain.Car, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Car), args.Error(1)
}

func (m *CarRepository) GetAll(ctx context.Context, limit, offset int) ([]domain.Car, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]domain.Car), args.Error(1)
}

func (m *CarRepository) Update(ctx context.Context, id string, car domain.Car) (*domain.Car, error) {
	args := m.Called(ctx, id, car)
	return args.Get(0).(*domain.Car), args.Error(1)
}

func (m *CarRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
