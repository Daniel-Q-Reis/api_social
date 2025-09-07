package v1

import (
	"encoding/json"
	"net/http"
	"strconv"

	"social/api/internal/controller/http/middleware"
	"social/api/internal/entity"
	"social/api/internal/repo"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type updateProfileRequest struct {
	Name     *string `json:"name,omitempty"`
	Bio      *string `json:"bio,omitempty"`
	ImageURL *string `json:"image_url,omitempty"`
}

type UserProfile struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Username  string  `json:"username"`
	Email     string  `json:"email"`
	Bio       *string `json:"bio,omitempty"`
	ImageURL  *string `json:"image_url,omitempty"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}

type usersResponse struct {
	Users []UserProfile `json:"users"`
}

func (h *Handler) getProfile(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	if username == "" {
		http.Error(w, "username is required", http.StatusBadRequest)
		return
	}

	user, err := h.userUseCase.GetProfile(r.Context(), username)
	if err != nil {
		if err == repo.ErrNotFound {
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}
		http.Error(w, "failed to get user profile", http.StatusInternalServerError)
		return
	}

	response := User{
		ID:        user.ID.String(),
		Name:      user.Name,
		Username:  user.Username,
		Email:     user.Email,
		Bio:       user.Bio,
		ImageURL:  user.ImageURL,
		CreatedAt: user.CreatedAt.String(),
		UpdatedAt: user.UpdatedAt.String(),
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		// Log the error but don't fail the request
		// In a production environment, you might want to log this
		_ = err // Explicitly ignore the error
	}
}

func (h *Handler) getMyProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserContextKey).(uuid.UUID)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := h.userUseCase.GetProfile(r.Context(), userID.String())
	if err != nil {
		if err == repo.ErrNotFound {
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}
		http.Error(w, "failed to get user profile", http.StatusInternalServerError)
		return
	}

	response := User{
		ID:        user.ID.String(),
		Name:      user.Name,
		Username:  user.Username,
		Email:     user.Email,
		Bio:       user.Bio,
		ImageURL:  user.ImageURL,
		CreatedAt: user.CreatedAt.String(),
		UpdatedAt: user.UpdatedAt.String(),
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		// Log the error but don't fail the request
		// In a production environment, you might want to log this
		_ = err // Explicitly ignore the error
	}
}

func (h *Handler) updateProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserContextKey).(uuid.UUID)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req updateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	user, err := h.userUseCase.UpdateProfile(r.Context(), userID, req.Name, req.Bio, req.ImageURL)
	if err != nil {
		if err == repo.ErrNotFound {
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := User{
		ID:        user.ID.String(),
		Name:      user.Name,
		Username:  user.Username,
		Email:     user.Email,
		Bio:       user.Bio,
		ImageURL:  user.ImageURL,
		CreatedAt: user.CreatedAt.String(),
		UpdatedAt: user.UpdatedAt.String(),
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		// Log the error but don't fail the request
		// In a production environment, you might want to log this
		_ = err // Explicitly ignore the error
	}
}

func (h *Handler) searchUsers(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "query parameter 'q' is required", http.StatusBadRequest)
		return
	}

	// Parse pagination parameters
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit <= 0 {
		limit = 20 // Default limit
	}

	if limit > 100 {
		limit = 100 // Maximum limit
	}

	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))

	if err != nil || offset < 0 {
		offset = 0 // Default offset
	}

	users, err := h.userUseCase.SearchUsers(r.Context(), query, limit, offset)
	if err != nil {
		http.Error(w, "failed to search users", http.StatusInternalServerError)
		return
	}

	responseUsers := make([]UserProfile, len(users))
	for i, user := range users {
		responseUsers[i] = UserProfile{
			ID:        user.ID.String(),
			Name:      user.Name,
			Username:  user.Username,
			Email:     user.Email,
			Bio:       user.Bio,
			ImageURL:  user.ImageURL,
			CreatedAt: user.CreatedAt.String(),
			UpdatedAt: user.UpdatedAt.String(),
		}
	}

	response := usersResponse{
		Users: responseUsers,
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		// Log the error but don't fail the request
		// In a production environment, you might want to log this
		_ = err // Explicitly ignore the error
	}
}

func (h *Handler) followUser(w http.ResponseWriter, r *http.Request) {
	h.handleUserFollow(w, r, true)
}

func (h *Handler) unfollowUser(w http.ResponseWriter, r *http.Request) {
	h.handleUserFollow(w, r, false)
}

func (h *Handler) handleUserFollow(w http.ResponseWriter, r *http.Request, isFollow bool) {
	// Get the authenticated user ID
	followerID, ok := r.Context().Value(middleware.UserContextKey).(uuid.UUID)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// Get the username to follow/unfollow from the URL
	username := chi.URLParam(r, "username")
	if username == "" {
		http.Error(w, "username is required", http.StatusBadRequest)
		return
	}

	// Get the user to follow/unfollow
	user, err := h.userUseCase.GetProfile(r.Context(), username)
	if err != nil {
		if err == repo.ErrNotFound {
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}
		http.Error(w, "failed to get user", http.StatusInternalServerError)
		return
	}

	userID := user.ID

	// Follow/unfollow the user
	if isFollow {
		err = h.interactionUseCase.FollowUser(r.Context(), userID, followerID)
	} else {
		err = h.interactionUseCase.UnfollowUser(r.Context(), userID, followerID)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) getFollowers(w http.ResponseWriter, r *http.Request) {
	h.handleUserRelation(w, r, true)
}

func (h *Handler) getFollowing(w http.ResponseWriter, r *http.Request) {
	h.handleUserRelation(w, r, false)
}

func (h *Handler) handleUserRelation(w http.ResponseWriter, r *http.Request, isFollowers bool) {
	username := chi.URLParam(r, "username")
	if username == "" {
		http.Error(w, "username is required", http.StatusBadRequest)
		return
	}

	// Parse pagination parameters
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit <= 0 {
		limit = 20 // Default limit
	}

	if limit > 100 {
		limit = 100 // Maximum limit
	}

	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))

	if err != nil || offset < 0 {
		offset = 0 // Default offset
	}

	var users []entity.User
	if isFollowers {
		users, err = h.interactionUseCase.GetFollowers(r.Context(), username, limit, offset)
	} else {
		users, err = h.interactionUseCase.GetFollowing(r.Context(), username, limit, offset)
	}

	if err != nil {
		if err == repo.ErrNotFound {
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}

		action := "followers"
		if !isFollowers {
			action = "following"
		}

		http.Error(w, "failed to get "+action, http.StatusInternalServerError)
		return
	}

	responseUsers := make([]UserProfile, len(users))
	for i, user := range users {
		responseUsers[i] = UserProfile{
			ID:        user.ID.String(),
			Name:      user.Name,
			Username:  user.Username,
			Email:     user.Email,
			Bio:       user.Bio,
			ImageURL:  user.ImageURL,
			CreatedAt: user.CreatedAt.String(),
			UpdatedAt: user.UpdatedAt.String(),
		}
	}

	response := usersResponse{
		Users: responseUsers,
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		// Log the error but don't fail the request
		// In a production environment, you might want to log this
		_ = err // Explicitly ignore the error
	}
}
