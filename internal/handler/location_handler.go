package handler

import (
	"PilaiteProject/internal/service"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type LocationHandler struct {
	locationService service.LocationService
}

func NewLocationHandler(s *service.LocationService) *LocationHandler {
	return &LocationHandler{
		locationService: *s,
	}
}

type CreateLocationRequest struct {
	Address   string  `json:"address"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func (h *LocationHandler) CreateLocation(w http.ResponseWriter, r *http.Request) {
	var req CreateLocationRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	location, err := h.locationService.InsertLocation(r.Context(), req.Address, req.Latitude, req.Longitude)
	if err != nil {
		http.Error(w, "Failed to create location: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(location)
}

func (h *LocationHandler) GetLocationById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "Location ID is needed", http.StatusBadRequest)
		return
	}

	locationId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, "Invalid location ID format", http.StatusBadRequest)
		return
	}
	location, err := h.locationService.GetLocationById(r.Context(), locationId)
	if err != nil {
		http.Error(w, "Failed to get location: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(location)
}

func (h *LocationHandler) GetAllLocations(w http.ResponseWriter, r *http.Request) {
	locations, err := h.locationService.GetAllLocations(r.Context())
	if err != nil {
		http.Error(w, "Failed to get locations: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(locations)
}
