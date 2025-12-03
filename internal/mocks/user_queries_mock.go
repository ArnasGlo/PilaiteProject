package mocks

import (
	"PilaiteProject/internal/db"
	"context"
)

type MockUserQueries struct {
	GetUserByIDFunc    func(ctx context.Context, id int64) (db.User, error)
	GetUserByEmailFunc func(ctx context.Context, email string) (db.User, error)
	GetAllUsersFunc    func(ctx context.Context) ([]db.User, error)
	InsertUserFunc     func(ctx context.Context, arg db.InsertUserParams) (db.User, error)
}

func (m *MockUserQueries) GetUserByID(ctx context.Context, id int64) (db.User, error) {
	return m.GetUserByIDFunc(ctx, id)
}
func (m *MockUserQueries) GetUserByEmail(ctx context.Context, email string) (db.User, error) {
	return m.GetUserByEmailFunc(ctx, email)
}
func (m *MockUserQueries) GetAllUsers(ctx context.Context) ([]db.User, error) {
	return m.GetAllUsersFunc(ctx)
}
func (m *MockUserQueries) InsertUser(ctx context.Context, arg db.InsertUserParams) (db.User, error) {
	return m.InsertUserFunc(ctx, arg)
}
