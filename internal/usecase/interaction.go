package usecase

import (
	"context"
	"fmt"

	"social/api/internal/entity"
	"social/api/internal/repo"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
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
	return s.handleLikeInteraction(ctx, postID, userID, true)
}

func (s *interactionService) UnlikePost(ctx context.Context, postID, userID uuid.UUID) error {
	return s.handleLikeInteraction(ctx, postID, userID, false)
}

func (s *interactionService) FollowUser(ctx context.Context, userID, followerID uuid.UUID) error {
	return s.handleFollowInteraction(ctx, userID, followerID, true)
}

func (s *interactionService) UnfollowUser(ctx context.Context, userID, followerID uuid.UUID) error {
	return s.handleFollowInteraction(ctx, userID, followerID, false)
}

func (s *interactionService) handleLikeInteraction(ctx context.Context, postID, userID uuid.UUID, isLike bool) error {
	action := "liking"
	if !isLike {
		action = "unliking"
	}

	log.Info().Str("postID", postID.String()).Str("userID", userID.String()).Msg(action + " post")

	var err error

	if isLike {
		like := &entity.Like{
			UserID: userID,
			PostID: postID,
		}
		err = s.likeRepo.Create(ctx, like)
	} else {
		err = s.likeRepo.Delete(ctx, userID, postID)
	}

	if err != nil {
		action = "like"
		if !isLike {
			action = "unlike"
		}

		log.Error().Err(err).Str("postID", postID.String()).Str("userID", userID.String()).Msg("failed to " + action + " post")

		return fmt.Errorf("failed to %s post: %w", action, err)
	}

	action = "liked"
	if !isLike {
		action = "unliked"
	}

	log.Info().Str("postID", postID.String()).Str("userID", userID.String()).Msg("post " + action + " successfully")

	return nil
}

func (s *interactionService) handleFollowInteraction(ctx context.Context, userID, followerID uuid.UUID, isFollow bool) error {
	action := "following"
	if !isFollow {
		action = "unfollowing"
	}

	log.Info().Str("userID", userID.String()).Str("followerID", followerID.String()).Msg(action + " user")

	// Prevent users from following themselves
	if isFollow && userID == followerID {
		log.Warn().Str("userID", userID.String()).Msg("user attempted to follow themselves")
		return fmt.Errorf("you cannot follow yourself")
	}

	var err error

	if isFollow {
		follow := &entity.Follow{
			UserID:     userID,
			FollowerID: followerID,
		}
		err = s.followRepo.Create(ctx, follow)
	} else {
		err = s.followRepo.Delete(ctx, userID, followerID)
	}

	if err != nil {
		action = "follow"
		if !isFollow {
			action = "unfollow"
		}

		log.Error().Err(err).Str("userID", userID.String()).Str("followerID", followerID.String()).Msg("failed to " + action + " user")

		return fmt.Errorf("failed to %s user: %w", action, err)
	}

	action = "followed"
	if !isFollow {
		action = "unfollowed"
	}

	log.Info().Str("userID", userID.String()).Str("followerID", followerID.String()).Msg("user " + action + " successfully")

	return nil
}

func (s *interactionService) GetFollowers(ctx context.Context, username string, limit, offset int) ([]entity.User, error) {
	return s.getUsersByRelation(ctx, username, limit, offset, true)
}

func (s *interactionService) GetFollowing(ctx context.Context, username string, limit, offset int) ([]entity.User, error) {
	return s.getUsersByRelation(ctx, username, limit, offset, false)
}

func (s *interactionService) getUsersByRelation(ctx context.Context, username string, limit, offset int, isFollowers bool) ([]entity.User, error) {
	log.Info().Str("username", username).Int("limit", limit).Int("offset", offset).Msg("fetching users by relation")

	user, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		log.Error().Err(err).Str("username", username).Msg("failed to get user")
		return nil, err
	}

	var users []entity.User

	if isFollowers {
		users, err = s.followRepo.GetFollowers(ctx, user.ID, limit, offset)
		if err != nil {
			log.Error().Err(err).Str("userID", user.ID.String()).Msg("failed to get followers")
			return nil, fmt.Errorf("failed to get followers: %w", err)
		}

		log.Info().Str("userID", user.ID.String()).Int("follower_count", len(users)).Msg("followers fetched successfully")
	} else {
		users, err = s.followRepo.GetFollowing(ctx, user.ID, limit, offset)
		if err != nil {
			log.Error().Err(err).Str("userID", user.ID.String()).Msg("failed to get following")
			return nil, fmt.Errorf("failed to get following: %w", err)
		}

		log.Info().Str("userID", user.ID.String()).Int("following_count", len(users)).Msg("following fetched successfully")
	}

	// Clear passwords before returning
	for i := range users {
		users[i].Password = ""
	}

	return users, nil
}
