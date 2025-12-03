package service

import (
	"PilaiteProject/internal/db"
	"context"
	"fmt"
)

type LocationService struct {
	queries *db.Queries
}

func NewLocationService(queries *db.Queries) *LocationService {
	return &LocationService{queries: queries}
}

func (s *LocationService) GetLocationById(ctx context.Context, id int64) (*db.Location, error) {
	location, err := s.queries.GetLocationByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &location, nil
}

func (s *LocationService) GetAllLocations(ctx context.Context) ([]db.Location, error) {
	locations, err := s.queries.GetAllLocations(ctx)
	if err != nil {
		return nil, err
	}
	return locations, nil
}

func (s *LocationService) InsertLocation(ctx context.Context, address string, latitude, longitude float64) (*db.Location, error) {
	// Basic validation
	if address == "" {
		return nil, fmt.Errorf("address cannot be empty")
	}
	// Validate latitude range (-90 to 90)
	if latitude < -90 || latitude > 90 {
		return nil, fmt.Errorf("latitude must be between -90 and 90, got: %f", latitude)
	}
	// Validate longitude range (-180 to 180)
	if longitude < -180 || longitude > 180 {
		return nil, fmt.Errorf("longitude must be between -180 and 180, got: %f", longitude)
	}

	location, err := s.queries.InsertLocation(ctx, db.InsertLocationParams{
		Address:   address,
		Latitude:  latitude,
		Longitude: longitude,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to insert location: %w", err)
	}

	return &location, nil
}
