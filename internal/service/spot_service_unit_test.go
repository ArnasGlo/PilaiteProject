package service

import (
	"PilaiteProject/internal/db"
	"PilaiteProject/internal/mocks"
	"context"
	"errors"
	"testing"
)

func TestGetSpotsByCategory_Success(t *testing.T) {
	mock := &mocks.MockSpotQueries{
		GetSpotsByCategoryWithDetailsFunc: func(ctx context.Context, category db.SpotCategory) ([]db.GetSpotsByCategoryWithDetailsRow, error) {
			return []db.GetSpotsByCategoryWithDetailsRow{
				{
					ID:        1,
					Category:  category,
					Name:      "Beach",
					Address:   "123 Beach St",
					ImageUrl:  "https://example.com/beach.jpg",
					Latitude:  54.6872,
					Longitude: 25.2797,
				},
				{
					ID:        2,
					Category:  category,
					Name:      "Mountain",
					Address:   "456 Mountain Rd",
					ImageUrl:  "https://example.com/mountain.jpg",
					Latitude:  54.7000,
					Longitude: 25.3000,
				},
			}, nil
		},
	}

	svc := NewSpotService(mock)
	ctx := context.Background()

	category := db.SpotCategorySlaptosVietos

	spots, err := svc.GetSpotsByCategoryWithDetails(ctx, category)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(spots) != 2 {
		t.Fatalf("expected 2 spots, got %d", len(spots))
	}
	if spots[0].Name != "Beach" {
		t.Fatalf("expected first spot name 'Beach', got %s", spots[0].Name)
	}

	t.Log(spots)
}

func TestGetSpotsByCategory_InvalidCategory(t *testing.T) {
	mock := &mocks.MockSpotQueries{}

	svc := NewSpotService(mock)
	ctx := context.Background()

	invalidCategory := db.SpotCategory("bad")

	_, err := svc.GetPublicSpotsByCategoryWithDetails(ctx, invalidCategory)
	if err == nil {
		t.Fatal("expected error for invalid category, got nil")
	}

	if err.Error() != "invalid category: bad" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGetSpotsByCategory_DBError(t *testing.T) {
	mock := &mocks.MockSpotQueries{
		GetSpotsByCategoryWithDetailsFunc: func(ctx context.Context, category db.SpotCategory) ([]db.GetSpotsByCategoryWithDetailsRow, error) {
			return nil, errors.New("db error")
		},
	}

	svc := NewSpotService(mock)
	ctx := context.Background()

	category := db.SpotCategorySlaptosVietos

	_, err := svc.GetSpotsByCategoryWithDetails(ctx, category)
	if err == nil {
		t.Fatal("expected error from DB, got nil")
	}
}

func TestGetPublicSpotsByCategory_Success(t *testing.T) {
	mock := &mocks.MockSpotQueries{
		GetSpotsByCategoryWithDetailsFunc: func(ctx context.Context, category db.SpotCategory) ([]db.GetSpotsByCategoryWithDetailsRow, error) {
			return []db.GetSpotsByCategoryWithDetailsRow{
				{
					ID:        1,
					Name:      "Park",
					Category:  category,
					Address:   "789 Park Ave",
					ImageUrl:  "https://example.com/park.jpg",
					Latitude:  54.6500,
					Longitude: 25.2500,
				},
				{
					ID:        2,
					Name:      "Lake",
					Category:  category,
					Address:   "321 Lake Dr",
					ImageUrl:  "https://example.com/lake.jpg",
					Latitude:  54.6800,
					Longitude: 25.2800,
				},
			}, nil
		},
	}

	svc := NewSpotService(mock)
	ctx := context.Background()
	category := db.SpotCategoryGamta

	spots, err := svc.GetPublicSpotsByCategoryWithDetails(ctx, category)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(spots) != 2 {
		t.Fatalf("expected 2 spots, got %d", len(spots))
	}
	t.Log(spots)
}

func TestGetPublicSpotsByCategory_SecretCategory(t *testing.T) {
	mock := &mocks.MockSpotQueries{
		GetPublicSpotsByCategoryWithDetailsFunc: func(ctx context.Context, category db.SpotCategory) ([]db.GetPublicSpotsByCategoryWithDetailsRow, error) {
			return nil, nil
		},
	}

	svc := NewSpotService(mock)
	ctx := context.Background()

	category := db.SpotCategorySlaptosVietos

	_, err := svc.GetPublicSpotsByCategoryWithDetails(ctx, category)
	if err == nil {
		t.Fatal("expected error for secret category, got nil")
	}

	expected := "secret category not accessible through public endpoint"
	if err.Error() != expected {
		t.Fatalf("unexpected error message: %v", err)
	}
}
