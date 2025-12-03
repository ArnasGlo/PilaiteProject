package interfaces

import (
	"context"
)

type SessionProvider interface {
	GetString(ctx context.Context, key string) string
	GetInt(ctx context.Context, key string) int
}
