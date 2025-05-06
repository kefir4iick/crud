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
