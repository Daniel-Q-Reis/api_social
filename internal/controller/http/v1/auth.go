package v1

import (
	"encoding/json"
	"net/http"

	"social/api/internal/repo"
)

type registerRequest struct {
	Name     string `json:"name" validate:"required"`
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type loginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type authResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type User struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Username  string  `json:"username"`
	Email     string  `json:"email"`
	Bio       *string `json:"bio,omitempty"`
	ImageURL  *string `json:"image_url,omitempty"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}

func (h *Handler) register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	user, err := h.userUseCase.Register(r.Context(), req.Name, req.Username, req.Email, req.Password)
	if err != nil {
		if err == repo.ErrDuplicateEmail {
			http.Error(w, "user with this email already exists", http.StatusConflict)
			return
		}
		if err == repo.ErrDuplicateUsername {
			http.Error(w, "user with this username already exists", http.StatusConflict)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// In a real implementation, you would generate a real JWT token here
	token, err := h.userUseCase.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		http.Error(w, "failed to generate token", http.StatusInternalServerError)
		return
	}

	response := authResponse{
		Token: token,
		User: User{
			ID:        user.ID.String(),
			Name:      user.Name,
			Username:  user.Username,
			Email:     user.Email,
			Bio:       user.Bio,
			ImageURL:  user.ImageURL,
			CreatedAt: user.CreatedAt.String(),
			UpdatedAt: user.UpdatedAt.String(),
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	token, err := h.userUseCase.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		if err == repo.ErrInvalidCredentials {
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}
		http.Error(w, "failed to login", http.StatusInternalServerError)
		return
	}

	// In a real implementation, you would get the user details from the token
	// For now, we'll just return a placeholder
	response := authResponse{
		Token: token,
		User: User{
			ID:       "00000000-0000-0000-0000-000000000000",
			Name:     "John Doe",
			Username: "johndoe",
			Email:    req.Email,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}