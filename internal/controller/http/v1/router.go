package v1

import (
	"github.com/go-chi/chi/v5"
	"social/api/internal/controller/http/middleware"
	"social/api/internal/usecase"
)

type Handler struct {
	userUseCase        usecase.User
	postUseCase        usecase.Post
	commentUseCase     usecase.Comment
	interactionUseCase usecase.Interaction
}

func NewHandler(userUseCase usecase.User, postUseCase usecase.Post, commentUseCase usecase.Comment, interactionUseCase usecase.Interaction) *Handler {
	return &Handler{
		userUseCase:        userUseCase,
		postUseCase:        postUseCase,
		commentUseCase:     commentUseCase,
		interactionUseCase: interactionUseCase,
	}
}

func (h *Handler) RegisterRoutes(r *chi.Mux) {
	// Public routes
	r.Post("/register", h.register)
	r.Post("/login", h.login)

	// User routes
	r.Get("/users/{username}", h.getProfile)
	r.Get("/users/search", h.searchUsers)

	// Post routes
	r.Get("/posts/{postID}", h.getPostByID)
	r.Get("/users/{username}/posts", h.getPostsByUser)

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(middleware.Auth)

		// User routes
		r.Get("/profile", h.getMyProfile)
		r.Put("/profile", h.updateProfile)

		// Following routes
		r.Post("/users/{username}/follow", h.followUser)
		r.Delete("/users/{username}/follow", h.unfollowUser)
		r.Get("/users/{username}/followers", h.getFollowers)
		r.Get("/users/{username}/following", h.getFollowing)

		// Post routes
		r.Post("/posts", h.createPost)
		r.Put("/posts/{postID}", h.updatePost)
		r.Delete("/posts/{postID}", h.deletePost)
		r.Get("/feed", h.getFeed)

		// Like routes
		r.Post("/posts/{postID}/like", h.likePost)
		r.Delete("/posts/{postID}/like", h.unlikePost)

		// Comment routes
		r.Post("/posts/{postID}/comments", h.addComment)
		r.Get("/posts/{postID}/comments", h.getComments)
		r.Delete("/posts/{postID}/comments/{commentID}", h.deleteComment)
	})
}