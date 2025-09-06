package usecase

import (
	"context"

	"github.com/google/uuid"
	"social/api/internal/entity"
)

type User interface {
	Register(ctx context.Context, name, username, email, password string) (*entity.User, error)
	Login(ctx context.Context, email, password string) (string, error)
	GetProfile(ctx context.Context, username string) (*entity.User, error)
	UpdateProfile(ctx context.Context, userID uuid.UUID, name, bio *string, imageURL *string) (*entity.User, error)
	SearchUsers(ctx context.Context, query string, limit, offset int) ([]entity.User, error)
}

type Post interface {
	CreatePost(ctx context.Context, authorID uuid.UUID, content string, imageURL *string) (*entity.Post, error)
	GetPostByID(ctx context.Context, postID uuid.UUID) (*entity.Post, error)
	GetPostsByUser(ctx context.Context, username string, limit, offset int) ([]entity.Post, error)
	UpdatePost(ctx context.Context, postID, userID uuid.UUID, content string, imageURL *string) (*entity.Post, error)
	DeletePost(ctx context.Context, postID, userID uuid.UUID) error
	GetFeed(ctx context.Context, userID uuid.UUID, limit, offset int) ([]entity.Post, error)
}

type Comment interface {
	AddComment(ctx context.Context, postID, userID uuid.UUID, content string) (*entity.Comment, error)
	GetComments(ctx context.Context, postID uuid.UUID, limit, offset int) ([]entity.Comment, error)
	DeleteComment(ctx context.Context, commentID, userID uuid.UUID) error
}

type Interaction interface {
	LikePost(ctx context.Context, postID, userID uuid.UUID) error
	UnlikePost(ctx context.Context, postID, userID uuid.UUID) error
	FollowUser(ctx context.Context, userID, followerID uuid.UUID) error
	UnfollowUser(ctx context.Context, userID, followerID uuid.UUID) error
	GetFollowers(ctx context.Context, username string, limit, offset int) ([]entity.User, error)
	GetFollowing(ctx context.Context, username string, limit, offset int) ([]entity.User, error)
}