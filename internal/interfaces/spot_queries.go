package interfaces

import (
	"PilaiteProject/internal/db"
	"context"
)

type SpotQueries interface {
	InsertSpot(ctx context.Context, arg db.InsertSpotParams) (db.Spot, error)
	GetSpotByID(ctx context.Context, id int64) (db.Spot, error)
	GetAllSpots(ctx context.Context) ([]db.Spot, error)
	GetPublicSpotsWithDetails(ctx context.Context) ([]db.GetPublicSpotsWithDetailsRow, error)
	GetSpotsWithDetails(ctx context.Context) ([]db.GetSpotsWithDetailsRow, error)
	GetPublicSpotsByCategoryWithDetails(ctx context.Context, category db.SpotCategory) ([]db.GetPublicSpotsByCategoryWithDetailsRow, error)
	GetSpotsByCategoryWithDetails(ctx context.Context, category db.SpotCategory) ([]db.GetSpotsByCategoryWithDetailsRow, error)
}
