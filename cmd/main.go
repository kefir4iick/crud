package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/kefir4iick/crud/internal/api"
	"github.com/kefir4iick/crud/internal/handler"
	"github.com/kefir4iick/crud/internal/repository/postgres"
	"github.com/kefir4iick/crud/internal/service"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	connStr := buildConnectionString()
	db, err := postgres.NewDB(connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	repo := postgres.NewPostgresCarRepository(db)

	carService := service.NewCarService(repo)

	carHandler := handler.NewCarHandler(carService)

	r := chi.NewRouter()
	r.Mount("/cars", api.NewCarRouter(carHandler))

	port := getEnv("PORT", "8080")
	log.Printf("Starting server on :%s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func buildConnectionString() string {
	return "user=" + getEnv("DB_USER", "postgres") +
		" dbname=" + getEnv("DB_NAME", "postgres") +
		" password=" + getEnv("DB_PASSWORD", "postgres") +
		" host=" + getEnv("DB_HOST", "localhost") +
		" port=" + getEnv("DB_PORT", "5432") +
		" sslmode=" + getEnv("DB_SSLMODE", "disable")
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
