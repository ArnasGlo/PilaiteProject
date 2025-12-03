package mocks

import (
	"PilaiteProject/internal/db"
	"context"
)

type MockSpotQueries struct {
	GetSpotByIDFunc                         func(ctx context.Context, id int64) (db.Spot, error)
	GetAllSpotsFunc                         func(ctx context.Context) ([]db.Spot, error)
	InsertSpotFunc                          func(ctx context.Context, arg db.InsertSpotParams) (db.Spot, error)
	GetPublicSpotsWithDetailsFunc           func(ctx context.Context) ([]db.GetPublicSpotsWithDetailsRow, error)
	GetSpotsWithDetailsFunc                 func(ctx context.Context) ([]db.GetSpotsWithDetailsRow, error)
	GetPublicSpotsByCategoryWithDetailsFunc func(ctx context.Context, category db.SpotCategory) ([]db.GetPublicSpotsByCategoryWithDetailsRow, error)
	GetSpotsByCategoryWithDetailsFunc       func(ctx context.Context, category db.SpotCategory) ([]db.GetSpotsByCategoryWithDetailsRow, error)
}

func (m MockSpotQueries) GetSpotByID(ctx context.Context, id int64) (db.Spot, error) {
	return m.GetSpotByIDFunc(ctx, id)
}

func (m MockSpotQueries) GetAllSpots(ctx context.Context) ([]db.Spot, error) {
	return m.GetAllSpotsFunc(ctx)
}

func (m MockSpotQueries) InsertSpot(ctx context.Context, arg db.InsertSpotParams) (db.Spot, error) {
	return m.InsertSpotFunc(ctx, arg)
}

func (m MockSpotQueries) GetPublicSpotsWithDetails(ctx context.Context) ([]db.GetPublicSpotsWithDetailsRow, error) {
	return m.GetPublicSpotsWithDetailsFunc(ctx)
}

func (m MockSpotQueries) GetSpotsWithDetails(ctx context.Context) ([]db.GetSpotsWithDetailsRow, error) {
	return m.GetSpotsWithDetailsFunc(ctx)
}

func (m MockSpotQueries) GetPublicSpotsByCategoryWithDetails(ctx context.Context, category db.SpotCategory) ([]db.GetPublicSpotsByCategoryWithDetailsRow, error) {
	return m.GetPublicSpotsByCategoryWithDetailsFunc(ctx, category)
}

func (m MockSpotQueries) GetSpotsByCategoryWithDetails(ctx context.Context, category db.SpotCategory) ([]db.GetSpotsByCategoryWithDetailsRow, error) {
	return m.GetSpotsByCategoryWithDetailsFunc(ctx, category)
}
