package server

import (
	"PilaiteProject/internal/mocks"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequireAuth_Success(t *testing.T) {
	mockSession := &mocks.MockSessionManager{
		GetIntFunc: func(ctx context.Context, key string) int {
			if key == "userID" {
				return 42
			}
			return 0
		},
		GetStringFunc: func(ctx context.Context, key string) string {
			return ""
		},
	}

	middleware := NewAuthMiddleware(mockSession)

	called := false
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		// Assert context contains data
		if r.Context().Value(userIDKey).(int) != 42 {
			t.Fatalf("expected userID 42")
		}
	})

	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	middleware.RequireAuth(nextHandler).ServeHTTP(rr, req)

	if !called {
		t.Fatal("next handler was not called")
	}
}

func TestRequireAuth_Unauthorized(t *testing.T) {
	mockSession := &mocks.MockSessionManager{
		GetIntFunc: func(ctx context.Context, key string) int {
			return 0
		},
		GetStringFunc: func(ctx context.Context, key string) string {
			return ""
		},
	}

	middleware := NewAuthMiddleware(mockSession)

	called := false
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	})

	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	handler := middleware.RequireAuth(nextHandler)
	handler.ServeHTTP(rr, req)

	if called {
		t.Fatal("next handler SHOULD NOT have been called")
	}

	//if rr.Code != http.StatusUnauthorized {
	//	t.Fatalf("expected 401 Unauthorized, got %d", rr.Code)
	//}
}

func TestRequireAdmin_Success(t *testing.T) {
	mock := &mocks.MockSessionManager{
		GetIntFunc: func(ctx context.Context, key string) int {
			if key == "userID" {
				return 42
			}
			return 0
		},
		GetStringFunc: func(ctx context.Context, key string) string {
			if key == "role" {
				return "admin"
			}
			return ""
		},
	}

	middleware := NewAuthMiddleware(mock)

	called := false
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true

		// Assert correct context values
		if r.Context().Value(userIDKey).(int) != 42 {
			t.Fatalf("expected userID 42")
		}
		if r.Context().Value(userRoleKey).(string) != "admin" {
			t.Fatalf("expected role admin")
		}
	})

	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	handler := middleware.RequireAdmin(nextHandler)
	handler.ServeHTTP(rr, req)

	if !called {
		t.Fatal("next handler was NOT called")
	}

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", rr.Code)
	}
}
func TestRequireAdmin_Unauthorized(t *testing.T) {
	mock := &mocks.MockSessionManager{
		GetIntFunc: func(ctx context.Context, key string) int {
			return 0 // user not authenticated
		},
		GetStringFunc: func(ctx context.Context, key string) string {
			return ""
		},
	}

	middleware := NewAuthMiddleware(mock)

	called := false
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	})

	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	handler := middleware.RequireAdmin(nextHandler)
	handler.ServeHTTP(rr, req)

	if called {
		t.Fatal("next handler SHOULD NOT have been called")
	}

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 Unauthorized, got %d", rr.Code)
	}
}

func TestRequireAdmin_Forbidden(t *testing.T) {
	mock := &mocks.MockSessionManager{
		GetIntFunc: func(ctx context.Context, key string) int {
			return 100 // logged in user
		},
		GetStringFunc: func(ctx context.Context, key string) string {
			if key == "role" {
				return "user" // not admin
			}
			return ""
		},
	}

	middleware := NewAuthMiddleware(mock)

	called := false
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	})

	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	handler := middleware.RequireAdmin(nextHandler)
	handler.ServeHTTP(rr, req)

	if called {
		t.Fatal("next handler SHOULD NOT have been called")
	}

	if rr.Code != http.StatusForbidden {
		t.Fatalf("expected 403 Forbidden, got %d", rr.Code)
	}
}

func TestRequireGuest_Success(t *testing.T) {
	mock := &mocks.MockSessionManager{
		GetIntFunc: func(ctx context.Context, key string) int {
			return 0
		},
		GetStringFunc: func(ctx context.Context, key string) string {
			return ""
		},
	}

	middleware := NewAuthMiddleware(mock)

	called := false
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true

	})
	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	handler := middleware.RequireGuest(nextHandler)
	handler.ServeHTTP(rr, req)

	if !called {
		t.Fatal("next handler was NOT called")
	}

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", rr.Code)
	}
}

func TestRequireGuest_Forbidden(t *testing.T) {
	mock := &mocks.MockSessionManager{
		GetIntFunc: func(ctx context.Context, key string) int {
			return 37
		},
		GetStringFunc: func(ctx context.Context, key string) string {
			return ""
		},
	}

	middleware := NewAuthMiddleware(mock)

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})

	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	handler := middleware.RequireGuest(nextHandler)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusForbidden {
		t.Fatalf("expected 403 Forbidden, got %d", rr.Code)
	}
}
