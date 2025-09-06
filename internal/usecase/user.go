package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
	"social/api/internal/entity"
	"social/api/internal/repo"
)

type userService struct {
	userRepo repo.User
}

func NewUserUseCase(userRepo repo.User) User {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) Register(ctx context.Context, name, username, email, password string) (*entity.User, error) {
	log.Info().Str("email", email).Str("username", username).Msg("registering new user")

	// Validate input
	user := &entity.User{
		Name:     name,
		Username: username,
		Email:    email,
		Password: password,
	}

	if err := user.Validate(); err != nil {
		log.Warn().Err(err).Str("email", email).Msg("user validation failed")
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Check if user already exists
	_, err := s.userRepo.GetByEmail(ctx, email)
	if err == nil {
		log.Warn().Str("email", email).Msg("user with this email already exists")
		return nil, repo.ErrDuplicateEmail
	}
	if !errors.Is(err, repo.ErrNotFound) {
		log.Error().Err(err).Str("email", email).Msg("failed to check email")
		return nil, fmt.Errorf("failed to check email: %w", err)
	}

	_, err = s.userRepo.GetByUsername(ctx, username)
	if err == nil {
		log.Warn().Str("username", username).Msg("user with this username already exists")
		return nil, repo.ErrDuplicateUsername
	}
	if !errors.Is(err, repo.ErrNotFound) {
		log.Error().Err(err).Str("username", username).Msg("failed to check username")
		return nil, fmt.Errorf("failed to check username: %w", err)
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error().Err(err).Msg("failed to hash password")
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user.Password = string(hashedPassword)

	err = s.userRepo.Create(ctx, user)
	if err != nil {
		log.Error().Err(err).Msg("failed to create user")
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	log.Info().Str("userID", user.ID.String()).Str("username", username).Msg("user registered successfully")

	// Clear password before returning
	user.Password = ""
	return user, nil
}

func (s *userService) Login(ctx context.Context, email, password string) (string, error) {
	log.Info().Str("email", email).Msg("user login attempt")

	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			log.Warn().Str("email", email).Msg("login failed: user not found")
			return "", repo.ErrInvalidCredentials
		}
		log.Error().Err(err).Str("email", email).Msg("failed to get user")
		return "", fmt.Errorf("failed to get user: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		log.Warn().Str("email", email).Msg("login failed: invalid password")
		return "", repo.ErrInvalidCredentials
	}

	log.Info().Str("userID", user.ID.String()).Str("email", email).Msg("user logged in successfully")

	// In a real implementation, you would generate a JWT token here
	// For now, we'll just return a placeholder
	return "jwt-token-placeholder", nil
}

func (s *userService) GetProfile(ctx context.Context, username string) (*entity.User, error) {
	log.Info().Str("username", username).Msg("fetching user profile")

	user, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			log.Warn().Str("username", username).Msg("user not found")
			return nil, err
		}
		log.Error().Err(err).Str("username", username).Msg("failed to get user")
		return nil, err
	}

	log.Info().Str("userID", user.ID.String()).Str("username", username).Msg("user profile fetched successfully")

	// Clear password before returning
	user.Password = ""
	return user, nil
}

func (s *userService) UpdateProfile(ctx context.Context, userID uuid.UUID, name, bio, imageURL *string) (*entity.User, error) {
	log.Info().Str("userID", userID.String()).Msg("updating user profile")

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			log.Warn().Str("userID", userID.String()).Msg("user not found")
			return nil, err
		}
		log.Error().Err(err).Str("userID", userID.String()).Msg("failed to get user")
		return nil, err
	}

	if name != nil {
		user.Name = *name
	}
	if bio != nil {
		user.Bio = bio
	}
	if imageURL != nil {
		user.ImageURL = imageURL
	}

	// Validate updated user
	if err := user.Validate(); err != nil {
		log.Warn().Err(err).Str("userID", userID.String()).Msg("user validation failed")
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	err = s.userRepo.Update(ctx, user)
	if err != nil {
		log.Error().Err(err).Str("userID", userID.String()).Msg("failed to update profile")
		return nil, fmt.Errorf("failed to update profile: %w", err)
	}

	log.Info().Str("userID", user.ID.String()).Msg("user profile updated successfully")

	// Clear password before returning
	user.Password = ""
	return user, nil
}

func (s *userService) SearchUsers(ctx context.Context, query string, limit, offset int) ([]entity.User, error) {
	log.Info().Str("query", query).Int("limit", limit).Int("offset", offset).Msg("searching users")

	users, err := s.userRepo.Search(ctx, query, limit, offset)
	if err != nil {
		log.Error().Err(err).Str("query", query).Msg("failed to search users")
		return nil, fmt.Errorf("failed to search users: %w", err)
	}

	log.Info().Str("query", query).Int("result_count", len(users)).Msg("users search completed")

	// Clear passwords before returning
	for i := range users {
		users[i].Password = ""
	}

	return users, nil
}