package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
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
	log.Info().Str("postID", postID.String()).Str("userID", userID.String()).Msg("adding comment to post")

	// Verify post exists
	_, err := s.postRepo.GetByID(ctx, postID)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			log.Warn().Str("postID", postID.String()).Msg("post not found when adding comment")
			return nil, err
		}
		log.Error().Err(err).Str("postID", postID.String()).Msg("failed to get post")
		return nil, err
	}

	comment := &entity.Comment{
		PostID:   postID,
		AuthorID: userID,
		Content:  content,
	}

	// Validate input
	if err := comment.Validate(); err != nil {
		log.Warn().Err(err).Str("postID", postID.String()).Str("userID", userID.String()).Msg("comment validation failed")
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	err = s.commentRepo.Create(ctx, comment)
	if err != nil {
		log.Error().Err(err).Str("postID", postID.String()).Str("userID", userID.String()).Msg("failed to add comment")
		return nil, fmt.Errorf("failed to add comment: %w", err)
	}

	log.Info().Str("commentID", comment.ID.String()).Str("postID", postID.String()).Msg("comment added successfully")

	return comment, nil
}

func (s *commentService) GetComments(ctx context.Context, postID uuid.UUID, limit, offset int) ([]entity.Comment, error) {
	log.Info().Str("postID", postID.String()).Int("limit", limit).Int("offset", offset).Msg("fetching comments for post")

	// Verify post exists
	_, err := s.postRepo.GetByID(ctx, postID)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			log.Warn().Str("postID", postID.String()).Msg("post not found when fetching comments")
			return nil, err
		}
		log.Error().Err(err).Str("postID", postID.String()).Msg("failed to get post")
		return nil, err
	}

	comments, err := s.commentRepo.GetByPostID(ctx, postID, limit, offset)
	if err != nil {
		log.Error().Err(err).Str("postID", postID.String()).Msg("failed to get comments")
		return nil, fmt.Errorf("failed to get comments: %w", err)
	}

	log.Info().Str("postID", postID.String()).Int("comment_count", len(comments)).Msg("comments fetched successfully")

	return comments, nil
}

func (s *commentService) DeleteComment(ctx context.Context, commentID, userID uuid.UUID) error {
	log.Info().Str("commentID", commentID.String()).Str("userID", userID.String()).Msg("deleting comment")

	// In a real implementation, you would fetch the comment to verify ownership
	// For now, we'll just delete it directly
	
	err := s.commentRepo.Delete(ctx, commentID)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			log.Warn().Str("commentID", commentID.String()).Msg("comment not found")
			return err
		}
		log.Error().Err(err).Str("commentID", commentID.String()).Msg("failed to delete comment")
		return fmt.Errorf("failed to delete comment: %w", err)
	}

	log.Info().Str("commentID", commentID.String()).Msg("comment deleted successfully")

	return nil
}