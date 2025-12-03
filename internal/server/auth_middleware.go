package server

import (
	"PilaiteProject/internal/interfaces"
	"context"
	"encoding/json"
	"net/http"
)

// AuthMiddleware holds the session manager dependency
type AuthMiddleware struct {
	sessionManager interfaces.SessionProvider
}

func NewAuthMiddleware(sessionManager interfaces.SessionProvider) *AuthMiddleware {
	return &AuthMiddleware{
		sessionManager: sessionManager,
	}
}

// ====== CONTEXT KEYS ======

type contextKey string

const (
	userIDKey    contextKey = "userID"
	userRoleKey  contextKey = "userRole"
	userEmailKey contextKey = "userEmail"
)

// ====== MIDDLEWARE FUNCTIONS ======

// RequireAuth checks if user is logged in
// Use this for routes that require any authenticated user
func (m *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if user is authenticated
		//fmt.Printf("REQUIRE AUTH DEBUG: Checking auth, session token: %v\n", m.sessionManager.Token(r.Context()))

		userID := m.sessionManager.GetInt(r.Context(), "userID")
		//fmt.Printf("REQUIRE AUTH DEBUG: UserID from session: %d\n", userID)
		if userID == 0 {
			respondWithError(w, http.StatusUnauthorized, "Authentication required")
			return
		}

		// Get additional user info from session
		role := m.sessionManager.GetString(r.Context(), "role")
		email := m.sessionManager.GetString(r.Context(), "email")

		// Add user info to request context for use in handlers
		ctx := r.Context()
		ctx = context.WithValue(ctx, userIDKey, userID)
		ctx = context.WithValue(ctx, userRoleKey, role)
		ctx = context.WithValue(ctx, userEmailKey, email)

		// Continue to next handler
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireAdmin checks if user is logged in AND is an admin
// Use this for admin-only routes
func (m *AuthMiddleware) RequireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if user is authenticated
		userID := m.sessionManager.GetInt(r.Context(), "userID")
		if userID == 0 {
			respondWithError(w, http.StatusUnauthorized, "Authentication required")
			return
		}

		// Check if user is admin
		role := m.sessionManager.GetString(r.Context(), "role")
		if role != "admin" {
			respondWithError(w, http.StatusForbidden, "Admin access required")
			return
		}

		// Get additional user info from session
		email := m.sessionManager.GetString(r.Context(), "email")

		// Add user info to request context
		ctx := r.Context()
		ctx = context.WithValue(ctx, userIDKey, userID)
		ctx = context.WithValue(ctx, userRoleKey, role)
		ctx = context.WithValue(ctx, userEmailKey, email)

		// Continue to next handler
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireGuest ensures user is NOT logged in
// Use this for routes like login/register pages that shouldn't be accessible when authenticated
func (m *AuthMiddleware) RequireGuest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if user is authenticated
		userID := m.sessionManager.GetInt(r.Context(), "userID")
		if userID != 0 {
			respondWithError(w, http.StatusForbidden, "Already authenticated")
			return
		}

		// Continue to next handler
		next.ServeHTTP(w, r)
	})
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	response, _ := json.Marshal(map[string]string{"error": message})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// RequireRole checks if user has a specific role
// Use this for custom role-based access control
//func (m *AuthMiddleware) RequireRole(allowedRoles ...db.UserRole) func(http.Handler) http.Handler {
//	return func(next http.Handler) http.Handler {
//		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//			// Check if user is authenticated
//			userID := m.sessionManager.GetInt(r.Context(), "userID")
//			if userID == 0 {
//				respondWithError(w, http.StatusUnauthorized, "Authentication required")
//				return
//			}
//
//			// Get user role
//			roleStr := m.sessionManager.GetString(r.Context(), "role")
//			userRole := db.UserRole(roleStr)
//
//			// Check if user's role is in allowed roles
//			allowed := false
//			for _, allowedRole := range allowedRoles {
//				if userRole == allowedRole {
//					allowed = true
//					break
//				}
//			}
//
//			if !allowed {
//				respondWithError(w, http.StatusForbidden, "Insufficient permissions")
//				return
//			}
//
//			// Get additional user info from session
//			email := m.sessionManager.GetString(r.Context(), "email")
//
//			// Add user info to request context
//			ctx := r.Context()
//			ctx = context.WithValue(ctx, userIDKey, int64(userID))
//			ctx = context.WithValue(ctx, userRoleKey, roleStr)
//			ctx = context.WithValue(ctx, userEmailKey, email)
//
//			// Continue to next handler
//			next.ServeHTTP(w, r.WithContext(ctx))
//		})
//	}
//}
//
//// ====== HELPER FUNCTIONS TO GET USER INFO FROM CONTEXT ======
//
//// GetUserIDFromContext retrieves user ID from request context
//func GetUserIDFromContext(ctx context.Context) (int64, bool) {
//	userID, ok := ctx.Value(userIDKey).(int64)
//	return userID, ok
//}
//
//// GetUserRoleFromContext retrieves user role from request context
//func GetUserRoleFromContext(ctx context.Context) (string, bool) {
//	role, ok := ctx.Value(userRoleKey).(string)
//	return role, ok
//}
//
//// GetUserEmailFromContext retrieves user email from request context
//func GetUserEmailFromContext(ctx context.Context) (string, bool) {
//	email, ok := ctx.Value(userEmailKey).(string)
//	return email, ok
//}

// ====== UTILITY FUNCTION ======
