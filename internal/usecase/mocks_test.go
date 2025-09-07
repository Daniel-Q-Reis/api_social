package usecase_test

import (
	"context"

	"social/api/internal/entity"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

// UserRepoMock is a mock implementation of repo.User interface
type UserRepoMock struct {
	mock.Mock
}

func (m *UserRepoMock) Create(ctx context.Context, user *entity.User) error {
	args := m.Called(ctx, user)
	// Set a dummy ID if not already set
	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}
	return args.Error(0)
}

func (m *UserRepoMock) GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	args := m.Called(ctx, id)
	user, ok := args.Get(0).(*entity.User)
	if !ok && args.Get(0) != nil {
		return nil, args.Error(1)
	}
	return user, args.Error(1)
}

func (m *UserRepoMock) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	args := m.Called(ctx, email)
	user, ok := args.Get(0).(*entity.User)
	if !ok && args.Get(0) != nil {
		return nil, args.Error(1)
	}
	return user, args.Error(1)
}

func (m *UserRepoMock) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	args := m.Called(ctx, username)
	user, ok := args.Get(0).(*entity.User)
	if !ok && args.Get(0) != nil {
		return nil, args.Error(1)
	}
	return user, args.Error(1)
}

func (m *UserRepoMock) Update(ctx context.Context, user *entity.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *UserRepoMock) Search(ctx context.Context, query string, limit, offset int) ([]entity.User, error) {
	args := m.Called(ctx, query, limit, offset)
	users, ok := args.Get(0).([]entity.User)
	if !ok && args.Get(0) != nil {
		return nil, args.Error(1)
	}
	return users, args.Error(1)
}

// PostRepoMock is a mock implementation of repo.Post interface
type PostRepoMock struct {
	mock.Mock
}

func (m *PostRepoMock) Create(ctx context.Context, post *entity.Post) error {
	args := m.Called(ctx, post)
	// Set a dummy ID if not already set
	if post.ID == uuid.Nil {
		post.ID = uuid.New()
	}
	return args.Error(0)
}

func (m *PostRepoMock) GetByID(ctx context.Context, id uuid.UUID) (*entity.Post, error) {
	args := m.Called(ctx, id)
	post, ok := args.Get(0).(*entity.Post)
	if !ok && args.Get(0) != nil {
		return nil, args.Error(1)
	}
	return post, args.Error(1)
}

func (m *PostRepoMock) GetByAuthorID(ctx context.Context, authorID uuid.UUID, limit, offset int) ([]entity.Post, error) {
	args := m.Called(ctx, authorID, limit, offset)
	posts, ok := args.Get(0).([]entity.Post)
	if !ok && args.Get(0) != nil {
		return nil, args.Error(1)
	}
	return posts, args.Error(1)
}

func (m *PostRepoMock) GetFeed(ctx context.Context, userID uuid.UUID, limit, offset int) ([]entity.Post, error) {
	args := m.Called(ctx, userID, limit, offset)
	posts, ok := args.Get(0).([]entity.Post)
	if !ok && args.Get(0) != nil {
		return nil, args.Error(1)
	}
	return posts, args.Error(1)
}

func (m *PostRepoMock) Update(ctx context.Context, post *entity.Post) error {
	args := m.Called(ctx, post)
	return args.Error(0)
}

func (m *PostRepoMock) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// CommentRepoMock is a mock implementation of repo.Comment interface
type CommentRepoMock struct {
	mock.Mock
}

func (m *CommentRepoMock) Create(ctx context.Context, comment *entity.Comment) error {
	args := m.Called(ctx, comment)
	// Set a dummy ID if not already set
	if comment.ID == uuid.Nil {
		comment.ID = uuid.New()
	}
	return args.Error(0)
}

func (m *CommentRepoMock) GetByPostID(ctx context.Context, postID uuid.UUID, limit, offset int) ([]entity.Comment, error) {
	args := m.Called(ctx, postID, limit, offset)
	comments, ok := args.Get(0).([]entity.Comment)
	if !ok && args.Get(0) != nil {
		return nil, args.Error(1)
	}
	return comments, args.Error(1)
}

func (m *CommentRepoMock) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// LikeRepoMock is a mock implementation of repo.Like interface
type LikeRepoMock struct {
	mock.Mock
}

func (m *LikeRepoMock) Create(ctx context.Context, like *entity.Like) error {
	args := m.Called(ctx, like)
	return args.Error(0)
}

func (m *LikeRepoMock) Delete(ctx context.Context, userID, postID uuid.UUID) error {
	args := m.Called(ctx, userID, postID)
	return args.Error(0)
}

func (m *LikeRepoMock) Exists(ctx context.Context, userID, postID uuid.UUID) (bool, error) {
	args := m.Called(ctx, userID, postID)
	return args.Bool(0), args.Error(1)
}

// FollowRepoMock is a mock implementation of repo.Follow interface
type FollowRepoMock struct {
	mock.Mock
}

func (m *FollowRepoMock) Create(ctx context.Context, follow *entity.Follow) error {
	args := m.Called(ctx, follow)
	return args.Error(0)
}

func (m *FollowRepoMock) Delete(ctx context.Context, userID, followerID uuid.UUID) error {
	args := m.Called(ctx, userID, followerID)
	return args.Error(0)
}

func (m *FollowRepoMock) Exists(ctx context.Context, userID, followerID uuid.UUID) (bool, error) {
	args := m.Called(ctx, userID, followerID)
	return args.Bool(0), args.Error(1)
}

func (m *FollowRepoMock) GetFollowers(ctx context.Context, userID uuid.UUID, limit, offset int) ([]entity.User, error) {
	args := m.Called(ctx, userID, limit, offset)
	users, ok := args.Get(0).([]entity.User)
	if !ok && args.Get(0) != nil {
		return nil, args.Error(1)
	}
	return users, args.Error(1)
}

func (m *FollowRepoMock) GetFollowing(ctx context.Context, followerID uuid.UUID, limit, offset int) ([]entity.User, error) {
	args := m.Called(ctx, followerID, limit, offset)
	users, ok := args.Get(0).([]entity.User)
	if !ok && args.Get(0) != nil {
		return nil, args.Error(1)
	}
	return users, args.Error(1)
}
