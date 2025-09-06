package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
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
	log.Info().Str("postID", postID.String()).Str("userID", userID.String()).Msg("liking post")

	like := &entity.Like{
		UserID: userID,
		PostID: postID,
	}

	err := s.likeRepo.Create(ctx, like)
	if err != nil {
		log.Error().Err(err).Str("postID", postID.String()).Str("userID", userID.String()).Msg("failed to like post")
		return fmt.Errorf("failed to like post: %w", err)
	}

	log.Info().Str("postID", postID.String()).Str("userID", userID.String()).Msg("post liked successfully")

	return nil
}

func (s *interactionService) UnlikePost(ctx context.Context, postID, userID uuid.UUID) error {
	log.Info().Str("postID", postID.String()).Str("userID", userID.String()).Msg("unliking post")

	err := s.likeRepo.Delete(ctx, userID, postID)
	if err != nil {
		log.Error().Err(err).Str("postID", postID.String()).Str("userID", userID.String()).Msg("failed to unlike post")
		return fmt.Errorf("failed to unlike post: %w", err)
	}

	log.Info().Str("postID", postID.String()).Str("userID", userID.String()).Msg("post unliked successfully")

	return nil
}

func (s *interactionService) FollowUser(ctx context.Context, userID, followerID uuid.UUID) error {
	log.Info().Str("userID", userID.String()).Str("followerID", followerID.String()).Msg("following user")

	// Prevent users from following themselves
	if userID == followerID {
		log.Warn().Str("userID", userID.String()).Msg("user attempted to follow themselves")
		return fmt.Errorf("you cannot follow yourself")
	}

	follow := &entity.Follow{
		UserID:     userID,
		FollowerID: followerID,
	}

	err := s.followRepo.Create(ctx, follow)
	if err != nil {
		log.Error().Err(err).Str("userID", userID.String()).Str("followerID", followerID.String()).Msg("failed to follow user")
		return fmt.Errorf("failed to follow user: %w", err)
	}

	log.Info().Str("userID", userID.String()).Str("followerID", followerID.String()).Msg("user followed successfully")

	return nil
}

func (s *interactionService) UnfollowUser(ctx context.Context, userID, followerID uuid.UUID) error {
	log.Info().Str("userID", userID.String()).Str("followerID", followerID.String()).Msg("unfollowing user")

	err := s.followRepo.Delete(ctx, userID, followerID)
	if err != nil {
		log.Error().Err(err).Str("userID", userID.String()).Str("followerID", followerID.String()).Msg("failed to unfollow user")
		return fmt.Errorf("failed to unfollow user: %w", err)
	}

	log.Info().Str("userID", userID.String()).Str("followerID", followerID.String()).Msg("user unfollowed successfully")

	return nil
}

func (s *interactionService) GetFollowers(ctx context.Context, username string, limit, offset int) ([]entity.User, error) {
	log.Info().Str("username", username).Int("limit", limit).Int("offset", offset).Msg("fetching followers")

	user, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		log.Error().Err(err).Str("username", username).Msg("failed to get user")
		return nil, err
	}

	followers, err := s.followRepo.GetFollowers(ctx, user.ID, limit, offset)
	if err != nil {
		log.Error().Err(err).Str("userID", user.ID.String()).Msg("failed to get followers")
		return nil, fmt.Errorf("failed to get followers: %w", err)
	}

	log.Info().Str("userID", user.ID.String()).Int("follower_count", len(followers)).Msg("followers fetched successfully")

	// Clear passwords before returning
	for i := range followers {
		followers[i].Password = ""
	}

	return followers, nil
}

func (s *interactionService) GetFollowing(ctx context.Context, username string, limit, offset int) ([]entity.User, error) {
	log.Info().Str("username", username).Int("limit", limit).Int("offset", offset).Msg("fetching following")

	user, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		log.Error().Err(err).Str("username", username).Msg("failed to get user")
		return nil, err
	}

	following, err := s.followRepo.GetFollowing(ctx, user.ID, limit, offset)
	if err != nil {
		log.Error().Err(err).Str("userID", user.ID.String()).Msg("failed to get following")
		return nil, fmt.Errorf("failed to get following: %w", err)
	}

	log.Info().Str("userID", user.ID.String()).Int("following_count", len(following)).Msg("following fetched successfully")

	// Clear passwords before returning
	for i := range following {
		following[i].Password = ""
	}

	return following, nil
}