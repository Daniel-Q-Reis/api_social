package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"social/api/internal/entity"
	"social/api/internal/repo"
)

type interactionService struct {
	likeRepo   repo.Like
	followRepo repo.Follow
	userRepo   repo.User
}

func NewInteractionUseCase(likeRepo repo.Like, followRepo repo.Follow, userRepo repo.User) Interaction {
	return &interactionService{
		likeRepo:   likeRepo,
		followRepo: followRepo,
		userRepo:   userRepo,
	}
}

func (s *interactionService) LikePost(ctx context.Context, postID, userID uuid.UUID) error {
	like := &entity.Like{
		UserID: userID,
		PostID: postID,
	}

	err := s.likeRepo.Create(ctx, like)
	if err != nil {
		return fmt.Errorf("failed to like post: %w", err)
	}

	return nil
}

func (s *interactionService) UnlikePost(ctx context.Context, postID, userID uuid.UUID) error {
	err := s.likeRepo.Delete(ctx, userID, postID)
	if err != nil {
		return fmt.Errorf("failed to unlike post: %w", err)
	}

	return nil
}

func (s *interactionService) FollowUser(ctx context.Context, userID, followerID uuid.UUID) error {
	// Prevent users from following themselves
	if userID == followerID {
		return fmt.Errorf("you cannot follow yourself")
	}

	follow := &entity.Follow{
		UserID:     userID,
		FollowerID: followerID,
	}

	err := s.followRepo.Create(ctx, follow)
	if err != nil {
		return fmt.Errorf("failed to follow user: %w", err)
	}

	return nil
}

func (s *interactionService) UnfollowUser(ctx context.Context, userID, followerID uuid.UUID) error {
	err := s.followRepo.Delete(ctx, userID, followerID)
	if err != nil {
		return fmt.Errorf("failed to unfollow user: %w", err)
	}

	return nil
}

func (s *interactionService) GetFollowers(ctx context.Context, username string) ([]entity.User, error) {
	user, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	followers, err := s.followRepo.GetFollowers(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get followers: %w", err)
	}

	// Clear passwords before returning
	for i := range followers {
		followers[i].Password = ""
	}

	return followers, nil
}

func (s *interactionService) GetFollowing(ctx context.Context, username string) ([]entity.User, error) {
	user, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	following, err := s.followRepo.GetFollowing(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get following: %w", err)
	}

	// Clear passwords before returning
	for i := range following {
		following[i].Password = ""
	}

	return following, nil
}