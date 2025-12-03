package service

import (
	"PilaiteProject/internal/db"
	"PilaiteProject/internal/interfaces"
	"context"
	"fmt"
)

type UserService struct {
	queries interfaces.UserQueries
}

func NewUserService(queries interfaces.UserQueries) *UserService {
	return &UserService{queries: queries}
}

func (s *UserService) GetUserById(ctx context.Context, id int64) (*db.User, error) {
	user, err := s.queries.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*db.User, error) {
	user, err := s.queries.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserService) GetAllUsers(ctx context.Context) ([]db.User, error) {
	users, err := s.queries.GetAllUsers(ctx)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *UserService) InsertUser(ctx context.Context, email, password string, role db.UserRole) (*db.User, error) {
	// Basic validation
	if email == "" {
		return nil, fmt.Errorf("email cannot be empty")
	}
	if password == "" {
		return nil, fmt.Errorf("password cannot be empty")
	}
	if !role.Valid() {
		return nil, fmt.Errorf("invalid role")
	}

	user, err := s.queries.InsertUser(ctx, db.InsertUserParams{
		Email:    email,
		Password: password,
		Role:     role,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to insert user: %w", err)
	}

	return &user, nil
}

func (s *UserService) StringToUserRole(roleStr string) (db.UserRole, error) {
	role := db.UserRole(roleStr)
	if !role.Valid() {
		return "", fmt.Errorf("invalid user role: %s", s)
	}
	return role, nil
}
