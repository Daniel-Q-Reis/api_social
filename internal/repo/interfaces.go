package repo

import (
	"context"

	"social/api/internal/entity"

	"github.com/google/uuid"
)

type User interface {
	Create(ctx context.Context, user *entity.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	GetByUsername(ctx context.Context, username string) (*entity.User, error)
	Update(ctx context.Context, user *entity.User) error
	Search(ctx context.Context, query string, limit, offset int) ([]entity.User, error)
}

type Post interface {
	Create(ctx context.Context, post *entity.Post) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Post, error)
	GetByAuthorID(ctx context.Context, authorID uuid.UUID, limit, offset int) ([]entity.Post, error)
	GetFeed(ctx context.Context, userID uuid.UUID, limit, offset int) ([]entity.Post, error)
	Update(ctx context.Context, post *entity.Post) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type Comment interface {
	Create(ctx context.Context, comment *entity.Comment) error
	GetByPostID(ctx context.Context, postID uuid.UUID, limit, offset int) ([]entity.Comment, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type Like interface {
	Create(ctx context.Context, like *entity.Like) error
	Delete(ctx context.Context, userID, postID uuid.UUID) error
	Exists(ctx context.Context, userID, postID uuid.UUID) (bool, error)
}

type Follow interface {
	Create(ctx context.Context, follow *entity.Follow) error
	Delete(ctx context.Context, userID, followerID uuid.UUID) error
	Exists(ctx context.Context, userID, followerID uuid.UUID) (bool, error)
	GetFollowers(ctx context.Context, userID uuid.UUID, limit, offset int) ([]entity.User, error)
	GetFollowing(ctx context.Context, followerID uuid.UUID, limit, offset int) ([]entity.User, error)
}
