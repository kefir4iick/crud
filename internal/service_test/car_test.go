package service_test

import (
	"context"
	"testing"

	"github.com/kefir4iick/crud/internal/domain"
	"github.com/kefir4iick/crud/internal/service"
	"github.com/kefir4iick/crud/internal/service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateCar_Validation(t *testing.T) {
	tests := []struct {
		name    string
		input   domain.Car
		wantErr string
	}{
		{
			name:    "Empty make",
			input:   domain.Car{Make: "", Model: "Model", Year: 2020, Price: 10000},
			wantErr: "make is required",
		},
		{
			name:    "Make too long",
			input:   domain.Car{Make: string(make([]byte, 256)), Model: "Model", Year: 2020, Price: 10000},
			wantErr: "make must be less than 255 characters",
		},
		{
			name:    "Empty model",
			input:   domain.Car{Make: "Make", Model: "", Year: 2020, Price: 10000},
			wantErr: "model is required",
		},
		{
			name:    "Year too old",
			input:   domain.Car{Make: "Make", Model: "Model", Year: 1899, Price: 10000},
			wantErr: "year must be >= 1900",
		},
		{
			name:    "Negative price",
			input:   domain.Car{Make: "Make", Model: "Model", Year: 2020, Price: -1},
			wantErr: "price must be positive",
		},
		{
			name:    "Zero price",
			input:   domain.Car{Make: "Make", Model: "Model", Year: 2020, Price: 0},
			wantErr: "price must be positive",
		},
		{
			name:    "Valid input",
			input:   domain.Car{Make: "Toyota", Model: "Camry", Year: 2020, Price: 25000},
			wantErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(mocks.CarRepository)
			s := service.NewCarService(repo)

			if tt.wantErr == "" {
				repo.On("Create", mock.Anything, tt.input).Return(&tt.input, nil)
			}

			_, err := s.Create(context.Background(), tt.input)

			if tt.wantErr != "" {
				assert.ErrorContains(t, err, tt.wantErr)
				repo.AssertNotCalled(t, "Create")
			} else {
				assert.NoError(t, err)
				repo.AssertCalled(t, "Create", mock.Anything, tt.input)
			}
		})
	}
}

func TestGetCarByID(t *testing.T) {
	tests := []struct {
		name     string
		id       string
		mockCar  *domain.Car
		mockErr  error
		wantCar  *domain.Car
		wantErr  string
	}{
		{
			name:    "Success",
			id:      "1",
			mockCar: &domain.Car{ID: "1", Make: "Toyota", Model: "Camry", Year: 2020, Price: 25000},
			wantCar: &domain.Car{ID: "1", Make: "Toyota", Model: "Camry", Year: 2020, Price: 25000},
		},
		{
			name:    "Empty ID",
			id:      "",
			wantErr: "id is required",
		},
		{
			name:    "Not found",
			id:      "999",
			mockErr: domain.ErrCarNotFound,
			wantErr: "car not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(mocks.CarRepository)
			
			if tt.id != "" {
				repo.On("GetByID", mock.Anything, tt.id).Return(tt.mockCar, tt.mockErr)
			}

			s := service.NewCarService(repo)
			car, err := s.GetByID(context.Background(), tt.id)

			if tt.wantErr != "" {
				assert.ErrorContains(t, err, tt.wantErr)
				assert.Nil(t, car)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantCar, car)
			}
			
			if tt.id != "" {
				repo.AssertCalled(t, "GetByID", mock.Anything, tt.id)
			}
		})
	}
}

func TestUpdateCar_Validation(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		input   domain.UpdateCarInput
		mockCar *domain.Car
		mockErr error
		wantErr string
	}{
		{
			name: "Make too long",
			id:   "1",
			input: domain.UpdateCarInput{
				Make: stringPtr(string(make([]byte, 256))),
			},
			mockCar: &domain.Car{ID: "1", Make: "Toyota", Model: "Camry", Year: 2020, Price: 25000},
			wantErr: "make must be less than 255 characters",
		},
		{
			name: "Empty make",
			id:   "1",
			input: domain.UpdateCarInput{
				Make: stringPtr(""),
			},
			mockCar: &domain.Car{ID: "1", Make: "Toyota", Model: "Camry", Year: 2020, Price: 25000},
			wantErr: "make cannot be empty",
		},
		{
			name: "Empty model",
			id:   "1",
			input: domain.UpdateCarInput{
				Model: stringPtr(""),
			},
			mockCar: &domain.Car{ID: "1", Make: "Toyota", Model: "Camry", Year: 2020, Price: 25000},
			wantErr: "model cannot be empty",
		},
		{
			name: "Year too old",
			id:   "1",
			input: domain.UpdateCarInput{
				Year: intPtr(1899),
			},
			mockCar: &domain.Car{ID: "1", Make: "Toyota", Model: "Camry", Year: 2020, Price: 25000},
			wantErr: "year must be >= 1900",
		},
		{
			name: "Negative price",
			id:   "1",
			input: domain.UpdateCarInput{
				Price: intPtr(-1),
			},
			mockCar: &domain.Car{ID: "1", Make: "Toyota", Model: "Camry", Year: 2020, Price: 25000},
			wantErr: "price must be positive",
		},
		{
			name: "Zero price",
			id:   "1",
			input: domain.UpdateCarInput{
				Price: intPtr(0),
			},
			mockCar: &domain.Car{ID: "1", Make: "Toyota", Model: "Camry", Year: 2020, Price: 25000},
			wantErr: "price must be positive",
		},
		{
			name: "Valid partial update",
			id:   "1",
			input: domain.UpdateCarInput{
				Price: intPtr(30000),
			},
			mockCar: &domain.Car{ID: "1", Make: "Toyota", Model: "Camry", Year: 2020, Price: 25000},
		},
		{
			name: "Car not found",
			id:   "999",
			input: domain.UpdateCarInput{
				Price: intPtr(30000),
			},
			mockCar: nil,
			mockErr: domain.ErrCarNotFound,
			wantErr: "car not found",
		},
		{
			name: "Empty ID",
			id:   "",
			input: domain.UpdateCarInput{
				Price: intPtr(30000),
			},
			wantErr: "id is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(mocks.CarRepository)
			
			if tt.id != "" {
				if tt.mockErr != nil {
					repo.On("GetByID", mock.Anything, tt.id).Return(nil, tt.mockErr)
				} else if tt.mockCar != nil {
					repo.On("GetByID", mock.Anything, tt.id).Return(tt.mockCar, nil)
				}
			}
			
			if tt.name != "Car not found" && tt.id != "" && tt.mockCar != nil {
				updatedCar := *tt.mockCar
				if tt.input.Make != nil {
					updatedCar.Make = *tt.input.Make
				}
				if tt.input.Model != nil {
					updatedCar.Model = *tt.input.Model
				}
				if tt.input.Year != nil {
					updatedCar.Year = *tt.input.Year
				}
				if tt.input.Price != nil {
					updatedCar.Price = *tt.input.Price
				}
				repo.On("Update", mock.Anything, tt.id, updatedCar).Return(&updatedCar, nil)
			}

			s := service.NewCarService(repo)
			_, err := s.Update(context.Background(), tt.id, tt.input)

			if tt.wantErr != "" {
				assert.ErrorContains(t, err, tt.wantErr)
				if tt.id != "" && tt.wantErr != "id is required" {
					repo.AssertCalled(t, "GetByID", mock.Anything, tt.id)
				}
				repo.AssertNotCalled(t, "Update")
			} else {
				assert.NoError(t, err)
				if tt.id != "" {
					repo.AssertCalled(t, "GetByID", mock.Anything, tt.id)
					repo.AssertCalled(t, "Update", mock.Anything, tt.id, mock.Anything)
				}
			}
		})
	}
}

func stringPtr(s string) *string { return &s }
func intPtr(i int) *int         { return &i }
