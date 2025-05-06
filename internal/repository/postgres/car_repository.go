package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/kefir4iick/crud/internal/domain"
	_ "github.com/lib/pq"
)

type postgresCarRepository struct {
	db *sql.DB
}

func NewPostgresCarRepository(db *sql.DB) *postgresCarRepository {
	return &postgresCarRepository{db: db}
}

func (r *postgresCarRepository) Create(ctx context.Context, car domain.Car) (*domain.Car, error) {
	query := `
		INSERT INTO cars (id, make, model, year, price)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, make, model, year, price
	`

	err := r.db.QueryRowContext(ctx, query,
		car.ID,
		car.Make,
		car.Model,
		car.Year,
		car.Price,
	).Scan(
		&car.ID,
		&car.Make,
		&car.Model,
		&car.Year,
		&car.Price,
	)

	if err != nil {
		if err.Error() == "pq: duplicate key value violates unique constraint \"cars_pkey\"" {
			return nil, domain.ErrDuplicateCarID
		}
		return nil, fmt.Errorf("failed to create car: %w", err)
	}

	return &car, nil
}

func (r *postgresCarRepository) GetByID(ctx context.Context, id string) (*domain.Car, error) {
	query := `
		SELECT id, make, model, year, price
		FROM cars
		WHERE id = $1
	`

	var car domain.Car
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&car.ID,
		&car.Make,
		&car.Model,
		&car.Year,
		&car.Price,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrCarNotFound
		}
		return nil, fmt.Errorf("failed to get car: %w", err)
	}

	return &car, nil
}

func (r *postgresCarRepository) GetAll(ctx context.Context, limit, offset int) ([]domain.Car, error) {
	query := `
		SELECT id, make, model, year, price
		FROM cars
		ORDER BY id
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get cars: %w", err)
	}
	defer rows.Close()

	var cars []domain.Car
	for rows.Next() {
		var car domain.Car
		if err := rows.Scan(
			&car.ID,
			&car.Make,
			&car.Model,
			&car.Year,
			&car.Price,
		); err != nil {
			return nil, fmt.Errorf("failed to scan car: %w", err)
		}
		cars = append(cars, car)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return cars, nil
}

func (r *postgresCarRepository) Update(ctx context.Context, id string, car domain.Car) (*domain.Car, error) {
	query := `
		UPDATE cars
		SET make = $1, model = $2, year = $3, price = $4
		WHERE id = $5
		RETURNING id, make, model, year, price
	`

	var updatedCar domain.Car
	err := r.db.QueryRowContext(ctx, query,
		car.Make,
		car.Model,
		car.Year,
		car.Price,
		id,
	).Scan(
		&updatedCar.ID,
		&updatedCar.Make,
		&updatedCar.Model,
		&updatedCar.Year,
		&updatedCar.Price,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrCarNotFound
		}
		return nil, fmt.Errorf("failed to update car: %w", err)
	}

	return &updatedCar, nil
}

func (r *postgresCarRepository) Delete(ctx context.Context, id string) error {
	query := `
		DELETE FROM cars
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete car: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return domain.ErrCarNotFound
	}

	return nil
}
