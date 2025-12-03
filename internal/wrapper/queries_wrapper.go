package wrapper

import (
	"PilaiteProject/internal/db"
	"PilaiteProject/internal/interfaces"
)

type AppQueries struct {
	*db.Queries
}

// Verify that AppQueries implements the interfaces
var (
	_ interfaces.SpotQueries = (*AppQueries)(nil)
)

func NewAppQueries(q *db.Queries) *AppQueries {
	return &AppQueries{Queries: q}
}
