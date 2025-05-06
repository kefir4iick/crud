package handler

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/kefir4iick/crud/internal/domain"
	"github.com/kefir4iick/crud/internal/service"
)

type CarHandler struct {
	service service.CarService
}

func NewCarHandler(service service.CarService) *CarHandler {
	return &CarHandler{service: service}
}

func (h *CarHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input domain.Car
	if err := decodeJSON(r, &input); err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}

	car, err := h.service.Create(r.Context(), input)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	respondJSON(w, http.StatusCreated, car)
}

func (h *CarHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	car, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		respondError(w, http.StatusNotFound, err)
		return
	}

	respondJSON(w, http.StatusOK, car)
}

func (h *CarHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	cars, err := h.service.GetAll(r.Context(), limit, offset)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	respondJSON(w, http.StatusOK, cars)
}

func (h *CarHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var input domain.UpdateCarInput
	if err := decodeJSON(r, &input); err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}

	car, err := h.service.Update(r.Context(), id, input)
	if err != nil {
		respondError(w, http.StatusNotFound, err)
		return
	}

	respondJSON(w, http.StatusOK, car)
}

func (h *CarHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := h.service.Delete(r.Context(), id); err != nil {
		respondError(w, http.StatusNotFound, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
