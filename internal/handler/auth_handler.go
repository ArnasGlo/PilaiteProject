package handler

import (
	"PilaiteProject/internal/db"
	"PilaiteProject/internal/service"
	"encoding/json"
	"net/http"
	"unicode"

	"github.com/alexedwards/scs/v2"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	userService    *service.UserService
	sessionManager *scs.SessionManager
}

func NewAuthHandler(userService *service.UserService, sessionManager *scs.SessionManager) *AuthHandler {
	return &AuthHandler{
		userService:    userService,
		sessionManager: sessionManager,
	}
}

type RegisterRequest struct {
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Message string   `json:"message"`
	User    *UserDTO `json:"user,omitempty"`
}

type UserDTO struct {
	ID    int64       `json:"id"`
	Email string      `json:"email"`
	Role  db.UserRole `json:"role"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" || req.ConfirmPassword == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	if req.Password != req.ConfirmPassword {
		http.Error(w, "Passwords do not match", http.StatusBadRequest)
		return
	}

	if len(req.Password) < 8 {
		http.Error(w, "Password must be at least 8 characters", http.StatusBadRequest)
		return
	}

	if !hasUppercase(req.Password) {
		http.Error(w, "Password must have atleast one uppercase", http.StatusBadRequest)
		return
	}

	if !hasDigit(req.Password) {
		http.Error(w, "Password must have at least one digit", http.StatusBadRequest)
		return
	}

	existingUser, err := h.userService.GetUserByEmail(r.Context(), req.Email)
	if err != nil && existingUser != nil {
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to generate password", http.StatusInternalServerError)
		return
	}

	role, err := h.userService.StringToUserRole("user")
	if err != nil {
		http.Error(w, "Failed to parse user role", http.StatusInternalServerError)
		return
	}

	user, err := h.userService.InsertUser(
		r.Context(),
		req.Email,
		string(hashPassword),
		role,
	)
	if err != nil {
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(AuthResponse{
		Message: "Registration successful",
		User: &UserDTO{
			ID:    user.ID,
			Email: user.Email,
			Role:  user.Role,
		},
	})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	user, err := h.userService.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	err = h.sessionManager.RenewToken(r.Context())
	if err != nil {
		http.Error(w, "Session error", http.StatusInternalServerError)
		return
	}

	h.sessionManager.Put(r.Context(), "userID", int(user.ID))
	h.sessionManager.Put(r.Context(), "role", string(user.Role))
	h.sessionManager.Put(r.Context(), "email", user.Email)

	//testUserID := h.sessionManager.GetInt(r.Context(), "userID")
	//fmt.Printf("LOGIN DEBUG: Immediately after Put, GetInt returns: %d\n", testUserID)

	//fmt.Printf("LOGIN DEBUG: Session created for user %d\n", user.ID)
	//fmt.Printf("LOGIN DEBUG: USER EMAIL: %s\n", user.Email)
	//fmt.Printf("LOGIN DEBUG: Session token in context: %v\n", h.sessionManager.Token(r.Context()))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(AuthResponse{
		Message: "Login successful",
		User: &UserDTO{
			ID:    user.ID,
			Email: user.Email,
			Role:  user.Role,
		},
	})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	err := h.sessionManager.Destroy(r.Context())
	if err != nil {
		http.Error(w, "Error logging out", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(AuthResponse{Message: "Logout successful"})
}

func (h *AuthHandler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	userID := h.sessionManager.GetInt(r.Context(), "userID")
	if userID == 0 {
		http.Error(w, "Not authenticated", http.StatusUnauthorized)
		return
	}

	user, err := h.userService.GetUserById(r.Context(), int64(userID))
	if err != nil {
		http.Error(w, "Failed to get user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(UserDTO{
		ID:    user.ID,
		Email: user.Email,
		Role:  user.Role,
	})
}

func hasDigit(s string) bool {
	for _, char := range s {
		if unicode.IsDigit(char) {
			return true
		}
	}
	return false
}

func hasUppercase(s string) bool {
	for _, char := range s {
		if unicode.IsUpper(char) {
			return true
		}
	}
	return false
}
