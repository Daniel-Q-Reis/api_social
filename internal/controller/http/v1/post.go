package v1

import (
	"encoding/json"
	"net/http"
	"strconv"

	"social/api/internal/controller/http/middleware"
	"social/api/internal/repo"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type createPostRequest struct {
	Content  string  `json:"content" validate:"required"`
	ImageURL *string `json:"image_url,omitempty"`
}

type Post struct {
	ID        string  `json:"id"`
	AuthorID  string  `json:"author_id"`
	Content   string  `json:"content"`
	ImageURL  *string `json:"image_url,omitempty"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}

type postResponse struct {
	Post Post `json:"post"`
}

type postsResponse struct {
	Posts []Post `json:"posts"`
}

func (h *Handler) createPost(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserContextKey).(uuid.UUID)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req createPostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	post, err := h.postUseCase.CreatePost(r.Context(), userID, req.Content, req.ImageURL)
	if err != nil {
		// Check for validation errors
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := postResponse{
		Post: Post{
			ID:        post.ID.String(),
			AuthorID:  post.AuthorID.String(),
			Content:   post.Content,
			ImageURL:  post.ImageURL,
			CreatedAt: post.CreatedAt.String(),
			UpdatedAt: post.UpdatedAt.String(),
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		// Log the error but don't fail the request
		// In a production environment, you might want to log this
		_ = err // Explicitly ignore the error
	}
}

func (h *Handler) getPostByID(w http.ResponseWriter, r *http.Request) {
	postIDStr := chi.URLParam(r, "postID")
	if postIDStr == "" {
		http.Error(w, "post ID is required", http.StatusBadRequest)
		return
	}

	postID, err := uuid.Parse(postIDStr)
	if err != nil {
		http.Error(w, "invalid post ID", http.StatusBadRequest)
		return
	}

	post, err := h.postUseCase.GetPostByID(r.Context(), postID)
	if err != nil {
		if err == repo.ErrNotFound {
			http.Error(w, "post not found", http.StatusNotFound)
			return
		}
		http.Error(w, "failed to get post", http.StatusInternalServerError)
		return
	}

	response := postResponse{
		Post: Post{
			ID:        post.ID.String(),
			AuthorID:  post.AuthorID.String(),
			Content:   post.Content,
			ImageURL:  post.ImageURL,
			CreatedAt: post.CreatedAt.String(),
			UpdatedAt: post.UpdatedAt.String(),
		},
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		// Log the error but don't fail the request
		// In a production environment, you might want to log this
		_ = err // Explicitly ignore the error
	}
}

func (h *Handler) getPostsByUser(w http.ResponseWriter, r *http.Request) {
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

	posts, err := h.postUseCase.GetPostsByUser(r.Context(), username, limit, offset)
	if err != nil {
		if err == repo.ErrNotFound {
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}
		http.Error(w, "failed to get posts", http.StatusInternalServerError)
		return
	}

	responsePosts := make([]Post, len(posts))
	for i, post := range posts {
		responsePosts[i] = Post{
			ID:        post.ID.String(),
			AuthorID:  post.AuthorID.String(),
			Content:   post.Content,
			ImageURL:  post.ImageURL,
			CreatedAt: post.CreatedAt.String(),
			UpdatedAt: post.UpdatedAt.String(),
		}
	}

	response := postsResponse{
		Posts: responsePosts,
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		// Log the error but don't fail the request
		// In a production environment, you might want to log this
		_ = err // Explicitly ignore the error
	}
}

func (h *Handler) updatePost(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserContextKey).(uuid.UUID)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	postIDStr := chi.URLParam(r, "postID")
	if postIDStr == "" {
		http.Error(w, "post ID is required", http.StatusBadRequest)
		return
	}

	postID, err := uuid.Parse(postIDStr)
	if err != nil {
		http.Error(w, "invalid post ID", http.StatusBadRequest)
		return
	}

	var req createPostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	post, err := h.postUseCase.UpdatePost(r.Context(), postID, userID, req.Content, req.ImageURL)
	if err != nil {
		if err == repo.ErrNotFound {
			http.Error(w, "post not found", http.StatusNotFound)
			return
		}

		if err == repo.ErrUnauthorized {
			http.Error(w, "unauthorized", http.StatusForbidden)
			return
		}

		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := postResponse{
		Post: Post{
			ID:        post.ID.String(),
			AuthorID:  post.AuthorID.String(),
			Content:   post.Content,
			ImageURL:  post.ImageURL,
			CreatedAt: post.CreatedAt.String(),
			UpdatedAt: post.UpdatedAt.String(),
		},
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		// Log the error but don't fail the request
		// In a production environment, you might want to log this
		_ = err // Explicitly ignore the error
	}
}

func (h *Handler) deletePost(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserContextKey).(uuid.UUID)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	postIDStr := chi.URLParam(r, "postID")
	if postIDStr == "" {
		http.Error(w, "post ID is required", http.StatusBadRequest)
		return
	}

	postID, err := uuid.Parse(postIDStr)
	if err != nil {
		http.Error(w, "invalid post ID", http.StatusBadRequest)
		return
	}

	err = h.postUseCase.DeletePost(r.Context(), postID, userID)
	if err != nil {
		if err == repo.ErrNotFound {
			http.Error(w, "post not found", http.StatusNotFound)
			return
		}

		if err == repo.ErrUnauthorized {
			http.Error(w, "unauthorized", http.StatusForbidden)
			return
		}

		http.Error(w, "failed to delete post", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) getFeed(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserContextKey).(uuid.UUID)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
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

	posts, err := h.postUseCase.GetFeed(r.Context(), userID, limit, offset)
	if err != nil {
		http.Error(w, "failed to get feed", http.StatusInternalServerError)
		return
	}

	responsePosts := make([]Post, len(posts))
	for i, post := range posts {
		responsePosts[i] = Post{
			ID:        post.ID.String(),
			AuthorID:  post.AuthorID.String(),
			Content:   post.Content,
			ImageURL:  post.ImageURL,
			CreatedAt: post.CreatedAt.String(),
			UpdatedAt: post.UpdatedAt.String(),
		}
	}

	response := postsResponse{
		Posts: responsePosts,
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		// Log the error but don't fail the request
		// In a production environment, you might want to log this
		_ = err // Explicitly ignore the error
	}
}

func (h *Handler) likePost(w http.ResponseWriter, r *http.Request) {
	h.handlePostInteraction(w, r, true)
}

func (h *Handler) unlikePost(w http.ResponseWriter, r *http.Request) {
	h.handlePostInteraction(w, r, false)
}

func (h *Handler) handlePostInteraction(w http.ResponseWriter, r *http.Request, isLike bool) {
	userID, ok := r.Context().Value(middleware.UserContextKey).(uuid.UUID)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	postIDStr := chi.URLParam(r, "postID")
	if postIDStr == "" {
		http.Error(w, "post ID is required", http.StatusBadRequest)
		return
	}

	postID, err := uuid.Parse(postIDStr)
	if err != nil {
		http.Error(w, "invalid post ID", http.StatusBadRequest)
		return
	}

	if isLike {
		err = h.interactionUseCase.LikePost(r.Context(), postID, userID)
	} else {
		err = h.interactionUseCase.UnlikePost(r.Context(), postID, userID)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
