package interfaces

import (
	"PilaiteProject/internal/db"
	"context"
)

type UserQueries interface {
	GetUserByID(ctx context.Context, id int64) (db.User, error)
	GetUserByEmail(ctx context.Context, email string) (db.User, error)
	GetAllUsers(ctx context.Context) ([]db.User, error)
	InsertUser(ctx context.Context, arg db.InsertUserParams) (db.User, error)
}
