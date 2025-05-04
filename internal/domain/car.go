package domain

type Car struct {
	ID    string `json:"id"`
	Make  string `json:"make" validate:"required"`
	Model string `json:"model" validate:"required"`
	Year  int    `json:"year" validate:"gte=1900"`
	Price int    `json:"price" validate:"gt=0"`
}

type UpdateCarInput struct {
	Make  *string `json:"make"`
	Model *string `json:"model"`
	Year  *int    `json:"year"`
	Price *int    `json:"price"`
}
