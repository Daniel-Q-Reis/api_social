package v1

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"social/api/internal/controller/http/middleware"
)

type updateProfileRequest struct {
	Name     *string `json:"name,omitempty"`
	Bio      *string `json:"bio,omitempty"`
	ImageURL *string `json:"image_url,omitempty"`
}

type searchUsersResponse struct {
	Users []User `json:"users"`
}

func (h *Handler) getProfile(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	if username == "" {
		http.Error(w, "username is required", http.StatusBadRequest)
		return
	}

	user, err := h.userUseCase.GetProfile(r.Context(), username)
	if err != nil {
		http.Error(w, "user not found", http.StatusNotFound)
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
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) getMyProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserContextKey).(uuid.UUID)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := h.userUseCase.GetProfile(r.Context(), userID.String())
	if err != nil {
		http.Error(w, "user not found", http.StatusNotFound)
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
	json.NewEncoder(w).Encode(response)
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
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) searchUsers(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "query parameter 'q' is required", http.StatusBadRequest)
		return
	}

	users, err := h.userUseCase.SearchUsers(r.Context(), query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responseUsers := make([]User, len(users))
	for i, user := range users {
		responseUsers[i] = User{
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

	response := searchUsersResponse{
		Users: responseUsers,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) followUser(w http.ResponseWriter, r *http.Request) {
	// Get the authenticated user ID
	followerID, ok := r.Context().Value(middleware.UserContextKey).(uuid.UUID)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// Get the username to follow from the URL
	username := chi.URLParam(r, "username")
	if username == "" {
		http.Error(w, "username is required", http.StatusBadRequest)
		return
	}

	// Get the user to follow
	user, err := h.userUseCase.GetProfile(r.Context(), username)
	if err != nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	userID := user.ID

	// Follow the user
	err = h.interactionUseCase.FollowUser(r.Context(), userID, followerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) unfollowUser(w http.ResponseWriter, r *http.Request) {
	// Get the authenticated user ID
	followerID, ok := r.Context().Value(middleware.UserContextKey).(uuid.UUID)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// Get the username to unfollow from the URL
	username := chi.URLParam(r, "username")
	if username == "" {
		http.Error(w, "username is required", http.StatusBadRequest)
		return
	}

	// Get the user to unfollow
	user, err := h.userUseCase.GetProfile(r.Context(), username)
	if err != nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	userID := user.ID

	// Unfollow the user
	err = h.interactionUseCase.UnfollowUser(r.Context(), userID, followerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) getFollowers(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	if username == "" {
		http.Error(w, "username is required", http.StatusBadRequest)
		return
	}

	followers, err := h.interactionUseCase.GetFollowers(r.Context(), username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	responseUsers := make([]User, len(followers))
	for i, follower := range followers {
		responseUsers[i] = User{
			ID:        follower.ID.String(),
			Name:      follower.Name,
			Username:  follower.Username,
			Email:     follower.Email,
			Bio:       follower.Bio,
			ImageURL:  follower.ImageURL,
			CreatedAt: follower.CreatedAt.String(),
			UpdatedAt: follower.UpdatedAt.String(),
		}
	}

	response := searchUsersResponse{
		Users: responseUsers,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) getFollowing(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	if username == "" {
		http.Error(w, "username is required", http.StatusBadRequest)
		return
	}

	following, err := h.interactionUseCase.GetFollowing(r.Context(), username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	responseUsers := make([]User, len(following))
	for i, user := range following {
		responseUsers[i] = User{
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

	response := searchUsersResponse{
		Users: responseUsers,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}