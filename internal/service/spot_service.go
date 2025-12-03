package service

import (
	"PilaiteProject/internal/db"
	"PilaiteProject/internal/dto"
	"PilaiteProject/internal/interfaces"
	"context"
	"fmt"
)

type SpotService struct {
	queries interfaces.SpotQueries
}

func NewSpotService(queries interfaces.SpotQueries) *SpotService {
	return &SpotService{queries: queries}
}

func (s *SpotService) InsertSpot(ctx context.Context, category db.SpotCategory, name, description string, location_id int64) (*db.Spot, error) {
	// Basic validation
	if !category.Valid() {
		return nil, fmt.Errorf("Invalid category: %v", category)
	}
	if name == "" {
		return nil, fmt.Errorf("name cannot be empty")
	}
	if description == "" {
		return nil, fmt.Errorf("description cannot be empty")
	}
	if location_id <= 0 {
		return nil, fmt.Errorf("invalid location_id: %d", location_id)
	}

	spot, err := s.queries.InsertSpot(ctx, db.InsertSpotParams{
		Category:    category,
		Name:        name,
		Description: description,
		LocationID:  location_id,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to insert a spot: %w", err)
	}

	return &spot, nil
}

func (s *SpotService) GetSpotById(ctx context.Context, id int64) (*db.Spot, error) {
	spot, err := s.queries.GetSpotByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &spot, nil
}

func (s *SpotService) GetAllSpots(ctx context.Context) ([]db.Spot, error) {
	spots, err := s.queries.GetAllSpots(ctx)
	if err != nil {
		return nil, err
	}
	return spots, nil
}

func (s *SpotService) GetPublicSpotsWithDetails(ctx context.Context) ([]dto.SpotCardDTO, error) {
	rows, err := s.queries.GetPublicSpotsWithDetails(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get spots with details: %w", err)
	}

	dtos := make([]dto.SpotCardDTO, len(rows))
	for i, row := range rows {
		dtos[i] = dto.SpotCardDTO{
			ID:        row.ID,
			Name:      row.Name,
			Category:  string(row.Category),
			Address:   row.Address,
			ImageURL:  row.ImageUrl,
			Latitude:  row.Latitude,
			Longitude: row.Longitude,
		}
	}
	return dtos, nil
}

func (s SpotService) GetSpotsWithDetails(ctx context.Context) ([]dto.SpotCardDTO, error) {
	rows, err := s.queries.GetSpotsWithDetails(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get spots with details: %w", err)
	}

	dtos := make([]dto.SpotCardDTO, len(rows))
	for i, row := range rows {
		dtos[i] = dto.SpotCardDTO{
			ID:        row.ID,
			Name:      row.Name,
			Category:  string(row.Category),
			Address:   row.Address,
			ImageURL:  row.ImageUrl,
			Latitude:  row.Latitude,
			Longitude: row.Longitude,
		}
	}
	return dtos, nil
}

func (s *SpotService) GetPublicSpotsByCategoryWithDetails(ctx context.Context, category db.SpotCategory) ([]dto.SpotCardDTO, error) {
	if !category.Valid() {
		return nil, fmt.Errorf("invalid category: %v", category)
	}

	if category == db.SpotCategorySlaptosVietos {
		return nil, fmt.Errorf("secret category not accessible through public endpoint")
	}

	rows, err := s.queries.GetSpotsByCategoryWithDetails(ctx, category)
	if err != nil {
		return nil, fmt.Errorf("failed to get spots by category: %w", err)
	}

	dtos := make([]dto.SpotCardDTO, len(rows))
	for i, row := range rows {
		dtos[i] = dto.SpotCardDTO{
			ID:        row.ID,
			Name:      row.Name,
			Category:  string(row.Category),
			Address:   row.Address,
			ImageURL:  row.ImageUrl,
			Latitude:  row.Latitude,
			Longitude: row.Longitude,
		}
	}
	return dtos, nil
}

func (s *SpotService) GetSpotsByCategoryWithDetails(ctx context.Context, category db.SpotCategory) ([]dto.SpotCardDTO, error) {
	if !category.Valid() {
		return nil, fmt.Errorf("invalid category: %v", category)
	}

	rows, err := s.queries.GetSpotsByCategoryWithDetails(ctx, category)
	if err != nil {
		return nil, fmt.Errorf("failed to get spots by category: %w", err)
	}

	dtos := make([]dto.SpotCardDTO, len(rows))
	for i, row := range rows {
		dtos[i] = dto.SpotCardDTO{
			ID:        row.ID,
			Name:      row.Name,
			Category:  string(row.Category),
			Address:   row.Address,
			ImageURL:  row.ImageUrl,
			Latitude:  row.Latitude,
			Longitude: row.Longitude,
		}
	}
	return dtos, nil
}

func (s *SpotService) StringToSpotCategory(categoryStr string) (db.SpotCategory, error) {
	category := db.SpotCategory(categoryStr)
	if !category.Valid() {
		return "", fmt.Errorf("invalid spot category: %s", categoryStr)
	}
	return category, nil
}
