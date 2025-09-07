package usecase

import (
	"context"
	"errors"
	"fmt"

	"social/api/internal/entity"
	"social/api/internal/repo"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type postService struct {
	postRepo repo.Post
	userRepo repo.User
}

func NewPostUseCase(postRepo repo.Post, userRepo repo.User) Post {
	return &postService{
		postRepo: postRepo,
		userRepo: userRepo,
	}
}

func (s *postService) CreatePost(ctx context.Context, authorID uuid.UUID, content string, imageURL *string) (*entity.Post, error) {
	log.Info().Str("authorID", authorID.String()).Msg("creating new post")

	post := &entity.Post{
		AuthorID: authorID,
		Content:  content,
		ImageURL: imageURL,
	}

	// Validate input
	if err := post.Validate(); err != nil {
		log.Warn().Err(err).Str("authorID", authorID.String()).Msg("post validation failed")
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	err := s.postRepo.Create(ctx, post)
	if err != nil {
		log.Error().Err(err).Str("authorID", authorID.String()).Msg("failed to create post")
		return nil, fmt.Errorf("failed to create post: %w", err)
	}

	log.Info().Str("postID", post.ID.String()).Str("authorID", authorID.String()).Msg("post created successfully")

	return post, nil
}

func (s *postService) GetPostByID(ctx context.Context, postID uuid.UUID) (*entity.Post, error) {
	log.Info().Str("postID", postID.String()).Msg("fetching post by ID")

	post, err := s.postRepo.GetByID(ctx, postID)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			log.Warn().Str("postID", postID.String()).Msg("post not found")
			return nil, err
		}
		log.Error().Err(err).Str("postID", postID.String()).Msg("failed to get post")
		return nil, err
	}

	log.Info().Str("postID", post.ID.String()).Msg("post fetched successfully")

	return post, nil
}

func (s *postService) GetPostsByUser(ctx context.Context, username string, limit, offset int) ([]entity.Post, error) {
	log.Info().Str("username", username).Int("limit", limit).Int("offset", offset).Msg("fetching posts by user")

	user, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			log.Warn().Str("username", username).Msg("user not found")
			return nil, err
		}
		log.Error().Err(err).Str("username", username).Msg("failed to get user")
		return nil, err
	}

	posts, err := s.postRepo.GetByAuthorID(ctx, user.ID, limit, offset)
	if err != nil {
		log.Error().Err(err).Str("userID", user.ID.String()).Msg("failed to get posts")
		return nil, fmt.Errorf("failed to get posts: %w", err)
	}

	log.Info().Str("userID", user.ID.String()).Int("post_count", len(posts)).Msg("posts fetched successfully")

	return posts, nil
}

func (s *postService) UpdatePost(ctx context.Context, postID, userID uuid.UUID, content string, imageURL *string) (*entity.Post, error) {
	log.Info().Str("postID", postID.String()).Str("userID", userID.String()).Msg("updating post")

	post, err := s.postRepo.GetByID(ctx, postID)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			log.Warn().Str("postID", postID.String()).Msg("post not found")
			return nil, err
		}
		log.Error().Err(err).Str("postID", postID.String()).Msg("failed to get post")
		return nil, err
	}

	if post.AuthorID != userID {
		log.Warn().Str("postID", postID.String()).Str("userID", userID.String()).Msg("unauthorized post update attempt")
		return nil, repo.ErrUnauthorized
	}

	post.Content = content
	if imageURL != nil {
		post.ImageURL = imageURL
	}

	// Validate updated post
	if err := post.Validate(); err != nil {
		log.Warn().Err(err).Str("postID", postID.String()).Msg("post validation failed")
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	err = s.postRepo.Update(ctx, post)
	if err != nil {
		log.Error().Err(err).Str("postID", postID.String()).Msg("failed to update post")
		return nil, fmt.Errorf("failed to update post: %w", err)
	}

	log.Info().Str("postID", post.ID.String()).Msg("post updated successfully")

	return post, nil
}

func (s *postService) DeletePost(ctx context.Context, postID, userID uuid.UUID) error {
	log.Info().Str("postID", postID.String()).Str("userID", userID.String()).Msg("deleting post")

	post, err := s.postRepo.GetByID(ctx, postID)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			log.Warn().Str("postID", postID.String()).Msg("post not found")
			return err
		}
		log.Error().Err(err).Str("postID", postID.String()).Msg("failed to get post")
		return err
	}

	if post.AuthorID != userID {
		log.Warn().Str("postID", postID.String()).Str("userID", userID.String()).Msg("unauthorized post delete attempt")
		return repo.ErrUnauthorized
	}

	err = s.postRepo.Delete(ctx, postID)
	if err != nil {
		log.Error().Err(err).Str("postID", postID.String()).Msg("failed to delete post")
		return fmt.Errorf("failed to delete post: %w", err)
	}

	log.Info().Str("postID", postID.String()).Msg("post deleted successfully")

	return nil
}

func (s *postService) GetFeed(ctx context.Context, userID uuid.UUID, limit, offset int) ([]entity.Post, error) {
	log.Info().Str("userID", userID.String()).Int("limit", limit).Int("offset", offset).Msg("fetching user feed")

	posts, err := s.postRepo.GetFeed(ctx, userID, limit, offset)
	if err != nil {
		log.Error().Err(err).Str("userID", userID.String()).Msg("failed to get feed")
		return nil, fmt.Errorf("failed to get feed: %w", err)
	}

	log.Info().Str("userID", userID.String()).Int("post_count", len(posts)).Msg("feed fetched successfully")

	return posts, nil
}
