package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
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
	// Check if user already exists
	_, err := s.userRepo.GetByEmail(ctx, email)
	if err == nil {
		return nil, fmt.Errorf("user with this email already exists")
	}

	_, err = s.userRepo.GetByUsername(ctx, username)
	if err == nil {
		return nil, fmt.Errorf("user with this username already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &entity.User{
		Name:     name,
		Username: username,
		Email:    email,
		Password: string(hashedPassword),
	}

	err = s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Clear password before returning
	user.Password = ""
	return user, nil
}

func (s *userService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", fmt.Errorf("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", fmt.Errorf("invalid credentials")
	}

	// In a real implementation, you would generate a JWT token here
	// For now, we'll just return a placeholder
	return "jwt-token-placeholder", nil
}

func (s *userService) GetProfile(ctx context.Context, username string) (*entity.User, error) {
	user, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Clear password before returning
	user.Password = ""
	return user, nil
}

func (s *userService) UpdateProfile(ctx context.Context, userID uuid.UUID, name, bio, imageURL *string) (*entity.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
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

	err = s.userRepo.Update(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to update profile: %w", err)
	}

	// Clear password before returning
	user.Password = ""
	return user, nil
}

func (s *userService) SearchUsers(ctx context.Context, query string) ([]entity.User, error) {
	users, err := s.userRepo.Search(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to search users: %w", err)
	}

	// Clear passwords before returning
	for i := range users {
		users[i].Password = ""
	}

	return users, nil
}