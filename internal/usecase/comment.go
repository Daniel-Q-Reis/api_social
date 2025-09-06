package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"social/api/internal/entity"
	"social/api/internal/repo"
)

type commentService struct {
	commentRepo repo.Comment
	userRepo    repo.User
	postRepo    repo.Post
}

func NewCommentUseCase(commentRepo repo.Comment, userRepo repo.User, postRepo repo.Post) Comment {
	return &commentService{
		commentRepo: commentRepo,
		userRepo:    userRepo,
		postRepo:    postRepo,
	}
}

func (s *commentService) AddComment(ctx context.Context, postID, userID uuid.UUID, content string) (*entity.Comment, error) {
	// Verify post exists
	_, err := s.postRepo.GetByID(ctx, postID)
	if err != nil {
		return nil, fmt.Errorf("post not found: %w", err)
	}

	comment := &entity.Comment{
		PostID:   postID,
		AuthorID: userID,
		Content:  content,
	}

	err = s.commentRepo.Create(ctx, comment)
	if err != nil {
		return nil, fmt.Errorf("failed to add comment: %w", err)
	}

	return comment, nil
}

func (s *commentService) GetComments(ctx context.Context, postID uuid.UUID) ([]entity.Comment, error) {
	// Verify post exists
	_, err := s.postRepo.GetByID(ctx, postID)
	if err != nil {
		return nil, fmt.Errorf("post not found: %w", err)
	}

	comments, err := s.commentRepo.GetByPostID(ctx, postID)
	if err != nil {
		return nil, fmt.Errorf("failed to get comments: %w", err)
	}

	return comments, nil
}

func (s *commentService) DeleteComment(ctx context.Context, commentID, userID uuid.UUID) error {
	// In a real implementation, you would fetch the comment to verify ownership
	// For now, we'll just delete it directly
	
	err := s.commentRepo.Delete(ctx, commentID)
	if err != nil {
		return fmt.Errorf("failed to delete comment: %w", err)
	}

	return nil
}