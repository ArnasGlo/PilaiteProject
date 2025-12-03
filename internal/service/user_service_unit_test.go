package service

import (
	"PilaiteProject/internal/db"
	"PilaiteProject/internal/mocks"
	"context"
	"testing"
)

func TestInsertUser_EmptyEmail(t *testing.T) {
	s := NewUserService(&mocks.MockUserQueries{})

	_, err := s.InsertUser(context.Background(), "", "password", db.UserRoleUser)

	if err.Error() != "email cannot be empty" {
		t.Fatalf("expected email validation error: got %v", err)
	}
	t.Log("Email validation passed successfully!")
}

func TestInsertUser_EmptyPassword(t *testing.T) {
	s := NewUserService(&mocks.MockUserQueries{})

	_, err := s.InsertUser(context.Background(), "naujas", "", db.UserRoleUser)

	if err.Error() != "password cannot be empty" {
		t.Fatalf("expected password validation error: got %v", err)
	}
	t.Log("Password validation passed successfully!")
}

func TestInsertUser_InvalidRole(t *testing.T) {
	s := NewUserService(&mocks.MockUserQueries{})

	invalidRole := db.UserRole("HIM")
	_, err := s.InsertUser(context.Background(), "naujas", "password", invalidRole)
	if err.Error() != "invalid role" {
		t.Fatalf("expected invalid role: %v", err)
	}
	t.Log("Role validation passed successfully!")
}

func TestInsertUser_Success(t *testing.T) {
	mock := &mocks.MockUserQueries{}

	mock.InsertUserFunc = func(ctx context.Context, arg db.InsertUserParams) (db.User, error) {
		return db.User{
			ID:       1,
			Email:    arg.Email,
			Password: arg.Password,
			Role:     arg.Role,
		}, nil
	}

	s := NewUserService(mock)

	user, err := s.InsertUser(context.Background(), "test@example.com", "testPass", db.UserRoleUser)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.Email != "test@example.com" {
		t.Fatalf("expected email: test@example.com, got %s", user.Email)
	}
	if user.Role != db.UserRoleUser {
		t.Fatalf("expected role: admin, got %s", user.Role)
	}
	t.Logf("User is inserted successfully: ID=%d, Email=%s, Role=%s", user.ID, user.Email, user.Role)
}

func TestGetUserById_Success(t *testing.T) {
	mock := &mocks.MockUserQueries{}

	mock.GetUserByIDFunc = func(ctx context.Context, id int64) (db.User, error) {
		return db.User{ID: id, Email: "mock@example.com"}, nil
	}

	s := NewUserService(mock)

	expectedId := int64(10)

	user, err := s.GetUserById(context.Background(), 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.ID != expectedId {
		t.Fatalf("expected ID 10, got %d", user.ID)
	}
	t.Logf("User is retrieved successfully: ID=%d, email: %s ", user.ID, user.Email)
}

func TestGetUserByEmail_Success(t *testing.T) {
	mock := &mocks.MockUserQueries{}

	mock.GetUserByEmailFunc = func(ctx context.Context, email string) (db.User, error) {
		return db.User{ID: 1, Email: email}, nil
	}
	s := NewUserService(mock)

	expectedEmail := "kazkoks@tipas@gmail.com"

	user, err := s.GetUserByEmail(context.Background(), "kazkoks@tipas@gmail.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.Email != expectedEmail {
		t.Fatalf("expected to get user by email: kazkoks@tipas@gmail.com, got %s", user.Email)
	}
	t.Logf("User is retrieved successfully: ID=%d, email: %s ", user.ID, user.Email)
}

func TestGetAllUsers_Success(t *testing.T) {
	mock := &mocks.MockUserQueries{}

	mock.GetAllUsersFunc = func(ctx context.Context) ([]db.User, error) {
		return []db.User{
			{ID: 1, Email: "user1@example.com"},
			{ID: 2, Email: "user2@example.com"},
		}, nil
	}

	expectedEmail1 := "user1@example.com"
	expectedEmail2 := "user2@example.com"

	s := NewUserService(mock)

	users, err := s.GetAllUsers(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(users) != 2 {
		t.Fatalf("expected 2 users, got %d", len(users))
	}
	if users[0].Email != expectedEmail1 {
		t.Fatalf("unexpected first user email: %s", users[0].Email)
	}
	if users[1].Email != expectedEmail2 {
		t.Fatalf("unexpected second user email: %s", users[1].Email)
	}
	t.Logf("Got all users: %v", users)
}

func TestStringToUserRole_Invalid(t *testing.T) {
	s := NewUserService(nil)

	_, err := s.StringToUserRole("HIM")
	if err == nil {
		t.Fatal("expected error for invalid role")
	}
}
