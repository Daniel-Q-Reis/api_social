package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"social/api/internal/entity"
	"social/api/internal/repo"
)

type postService struct {
	postRepo   repo.Post
	userRepo   repo.User
}

func NewPostUseCase(postRepo repo.Post, userRepo repo.User) Post {
	return &postService{
		postRepo: postRepo,
		userRepo: userRepo,
	}
}

func (s *postService) CreatePost(ctx context.Context, authorID uuid.UUID, content string, imageURL *string) (*entity.Post, error) {
	post := &entity.Post{
		AuthorID: authorID,
		Content:  content,
		ImageURL: imageURL,
	}

	err := s.postRepo.Create(ctx, post)
	if err != nil {
		return nil, fmt.Errorf("failed to create post: %w", err)
	}

	return post, nil
}

func (s *postService) GetPostByID(ctx context.Context, postID uuid.UUID) (*entity.Post, error) {
	post, err := s.postRepo.GetByID(ctx, postID)
	if err != nil {
		return nil, fmt.Errorf("post not found: %w", err)
	}

	return post, nil
}

func (s *postService) GetPostsByUser(ctx context.Context, username string) ([]entity.Post, error) {
	user, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	posts, err := s.postRepo.GetByAuthorID(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get posts: %w", err)
	}

	return posts, nil
}

func (s *postService) UpdatePost(ctx context.Context, postID, userID uuid.UUID, content string, imageURL *string) (*entity.Post, error) {
	post, err := s.postRepo.GetByID(ctx, postID)
	if err != nil {
		return nil, fmt.Errorf("post not found: %w", err)
	}

	if post.AuthorID != userID {
		return nil, fmt.Errorf("unauthorized: you can only update your own posts")
	}

	post.Content = content
	if imageURL != nil {
		post.ImageURL = imageURL
	}

	err = s.postRepo.Update(ctx, post)
	if err != nil {
		return nil, fmt.Errorf("failed to update post: %w", err)
	}

	return post, nil
}

func (s *postService) DeletePost(ctx context.Context, postID, userID uuid.UUID) error {
	post, err := s.postRepo.GetByID(ctx, postID)
	if err != nil {
		return fmt.Errorf("post not found: %w", err)
	}

	if post.AuthorID != userID {
		return fmt.Errorf("unauthorized: you can only delete your own posts")
	}

	err = s.postRepo.Delete(ctx, postID)
	if err != nil {
		return fmt.Errorf("failed to delete post: %w", err)
	}

	return nil
}

func (s *postService) GetFeed(ctx context.Context, userID uuid.UUID) ([]entity.Post, error) {
	posts, err := s.postRepo.GetFeed(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get feed: %w", err)
	}

	return posts, nil
}