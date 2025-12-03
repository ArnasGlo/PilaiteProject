package handler

import (
	"PilaiteProject/internal/db"
	"PilaiteProject/internal/service"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type SpotHandler struct {
	spotService service.SpotService
}

func NewSpotHandler(s *service.SpotService) *SpotHandler {
	return &SpotHandler{
		spotService: *s,
	}
}

type CreateSpotRequest struct {
	Category    db.SpotCategory `json:"category"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	LocationID  int64           `json:"location_id"`
}

func (h *SpotHandler) InsertSpot(w http.ResponseWriter, r *http.Request) {
	var req CreateSpotRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	spot, err := h.spotService.InsertSpot(r.Context(), req.Category, req.Name, req.Description, req.LocationID)
	if err != nil {
		http.Error(w, "Failed to create spot: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(spot)
}

func (h *SpotHandler) GetSpotById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "Spot ID is needed", http.StatusBadRequest)
		return
	}

	spotId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, "Invalid spot ID format", http.StatusBadRequest)
		return
	}
	spot, err := h.spotService.GetSpotById(r.Context(), spotId)
	if err != nil {
		http.Error(w, "Failed to get spot: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(spot)
}

func (h *SpotHandler) GetPublicSpotsWithDetails(w http.ResponseWriter, r *http.Request) {
	spots, err := h.spotService.GetPublicSpotsWithDetails(r.Context())
	if err != nil {
		http.Error(w, "Failed to get spots: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(spots)
}

func (h *SpotHandler) GetSpotsWithDetails(w http.ResponseWriter, r *http.Request) {
	spots, err := h.spotService.GetSpotsWithDetails(r.Context())
	if err != nil {
		http.Error(w, "Failed to get spots: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(spots)
}

func (h *SpotHandler) GetPublicSpotsByCategoryWithDetails(w http.ResponseWriter, r *http.Request) {
	categoryStr := chi.URLParam(r, "category")
	if categoryStr == "" {
		http.Error(w, "Category is required", http.StatusBadRequest)
		return
	}

	category, err := h.spotService.StringToSpotCategory(categoryStr)
	if err != nil {
		http.Error(w, "Invalid category: "+err.Error(), http.StatusBadRequest)
		return
	}

	spots, err := h.spotService.GetPublicSpotsByCategoryWithDetails(r.Context(), category)
	if err != nil {
		http.Error(w, "Failed to get spots: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(spots)
}

func (h *SpotHandler) GetSpotsByCategoryWithDetails(w http.ResponseWriter, r *http.Request) {
	categoryStr := chi.URLParam(r, "category")
	if categoryStr == "" {
		http.Error(w, "Category is required", http.StatusBadRequest)
		return
	}

	category, err := h.spotService.StringToSpotCategory(categoryStr)
	if err != nil {
		http.Error(w, "Invalid category: "+err.Error(), http.StatusBadRequest)
		return
	}

	spots, err := h.spotService.GetSpotsByCategoryWithDetails(r.Context(), category)
	if err != nil {
		http.Error(w, "Failed to get spots by category: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(spots)
}

//func (h *SpotHandler) GetSecretSpotsByCategory(w http.ResponseWriter, r *http.Request) {
//	// This handler is only for secret category and is protected by auth middleware
//	category := db.SpotCategorySlaptosVietos
//
//	spots, err := h.spotService.GetSpotsByCategory(r.Context(), category)
//	if err != nil {
//		http.Error(w, "Failed to get spots: "+err.Error(), http.StatusInternalServerError)
//		return
//	}
//
//	w.Header().Set("Content-Type", "application/json")
//	json.NewEncoder(w).Encode(spots)
//}

//func (h *SpotHandler) GetSpotWithLocation(w http.ResponseWriter, r *http.Request) {
//	id := chi.URLParam(r, "id")
//	if id == "" {
//		http.Error(w, "Spot ID is needed", http.StatusBadRequest)
//		return
//	}
//
//	spotId, err := strconv.ParseInt(id, 10, 64)
//	if err != nil {
//		http.Error(w, "Invalid spot ID format", http.StatusBadRequest)
//		return
//	}
//
//	spotWithLocation, err := h.spotService.GetSpotWithLocation(r.Context(), spotId)
//	if err != nil {
//		http.Error(w, "Failed to get spot with location: "+err.Error(), http.StatusInternalServerError)
//		return
//	}
//
//	w.Header().Set("Content-Type", "application/json")
//	json.NewEncoder(w).Encode(spotWithLocation)
//}
