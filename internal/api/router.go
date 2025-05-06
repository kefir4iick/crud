package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/kefir4iick/crud/internal/handler"
)

func NewCarRouter(h *handler.CarHandler) chi.Router {
	r := chi.NewRouter()

	r.Post("/", h.Create)
	r.Get("/", h.GetAll)
	r.Get("/{id}", h.GetByID)
	r.Put("/{id}", h.Update)
	r.Patch("/{id}", h.Update) 
	r.Delete("/{id}", h.Delete)

	return r
}
