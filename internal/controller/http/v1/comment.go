package v1

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"social/api/internal/controller/http/middleware"
)

type addCommentRequest struct {
	Content string `json:"content" validate:"required"`
}

type Comment struct {
	ID        string `json:"id"`
	PostID    string `json:"post_id"`
	AuthorID  string `json:"author_id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}

type commentResponse struct {
	Comment Comment `json:"comment"`
}

type commentsResponse struct {
	Comments []Comment `json:"comments"`
}

func (h *Handler) addComment(w http.ResponseWriter, r *http.Request) {
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

	var req addCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	comment, err := h.commentUseCase.AddComment(r.Context(), postID, userID, req.Content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := commentResponse{
		Comment: Comment{
			ID:        comment.ID.String(),
			PostID:    comment.PostID.String(),
			AuthorID:  comment.AuthorID.String(),
			Content:   comment.Content,
			CreatedAt: comment.CreatedAt.String(),
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) getComments(w http.ResponseWriter, r *http.Request) {
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

	comments, err := h.commentUseCase.GetComments(r.Context(), postID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responseComments := make([]Comment, len(comments))
	for i, comment := range comments {
		responseComments[i] = Comment{
			ID:        comment.ID.String(),
			PostID:    comment.PostID.String(),
			AuthorID:  comment.AuthorID.String(),
			Content:   comment.Content,
			CreatedAt: comment.CreatedAt.String(),
		}
	}

	response := commentsResponse{
		Comments: responseComments,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) deleteComment(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserContextKey).(uuid.UUID)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	commentIDStr := chi.URLParam(r, "commentID")
	if commentIDStr == "" {
		http.Error(w, "comment ID is required", http.StatusBadRequest)
		return
	}

	commentID, err := uuid.Parse(commentIDStr)
	if err != nil {
		http.Error(w, "invalid comment ID", http.StatusBadRequest)
		return
	}

	err = h.commentUseCase.DeleteComment(r.Context(), commentID, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}