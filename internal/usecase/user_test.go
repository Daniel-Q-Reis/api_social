package usecase_test

import (
	"context"
	"errors"
	"testing"

	"social/api/internal/entity"
	"social/api/internal/repo"
	"social/api/internal/usecase"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestUserService_Register(t *testing.T) {
	// Test cases
	tests := []struct {
		name          string
		nameInput     string
		usernameInput string
		emailInput    string
		passwordInput string
		setupMock     func(*UserRepoMock)
		expectedUser  *entity.User
		expectedError error
	}{
		{
			name:          "Success",
			nameInput:     "John Doe",
			usernameInput: "johndoe",
			emailInput:    "john@example.com",
			passwordInput: "password123",
			setupMock: func(mockRepo *UserRepoMock) {
				// Mock GetByEmail to return ErrNotFound (user doesn't exist)
				mockRepo.On("GetByEmail", mock.Anything, "john@example.com").Return(nil, repo.ErrNotFound)
				// Mock GetByUsername to return ErrNotFound (user doesn't exist)
				mockRepo.On("GetByUsername", mock.Anything, "johndoe").Return(nil, repo.ErrNotFound)
				// Mock Create to succeed
				mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(user *entity.User) bool {
					// Verify that the password was hashed
					err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte("password123"))
					return err == nil && user.Name == "John Doe" && user.Username == "johndoe" && user.Email == "john@example.com"
				})).Return(nil)
			},
			expectedUser: &entity.User{
				Name:     "John Doe",
				Username: "johndoe",
				Email:    "john@example.com",
				Password: "", // Password should be cleared in the returned user
			},
			expectedError: nil,
		},
		{
			name:          "ValidationFailed_InvalidEmail",
			nameInput:     "John Doe",
			usernameInput: "johndoe",
			emailInput:    "invalid-email", // Invalid email
			passwordInput: "password123",
			setupMock:     func(mockRepo *UserRepoMock) {},
			expectedUser:  nil,
			expectedError: errors.New("validation failed"),
		},
		{
			name:          "ValidationFailed_ShortPassword",
			nameInput:     "John Doe",
			usernameInput: "johndoe",
			emailInput:    "john@example.com",
			passwordInput: "123", // Too short
			setupMock:     func(mockRepo *UserRepoMock) {},
			expectedUser:  nil,
			expectedError: errors.New("validation failed"),
		},
		{
			name:          "EmailAlreadyExists",
			nameInput:     "John Doe",
			usernameInput: "johndoe",
			emailInput:    "john@example.com",
			passwordInput: "password123",
			setupMock: func(mockRepo *UserRepoMock) {
				// Mock GetByEmail to return an existing user
				existingUser := &entity.User{Email: "john@example.com"}
				mockRepo.On("GetByEmail", mock.Anything, "john@example.com").Return(existingUser, nil)
			},
			expectedUser:  nil,
			expectedError: repo.ErrDuplicateEmail,
		},
		{
			name:          "UsernameAlreadyExists",
			nameInput:     "John Doe",
			usernameInput: "johndoe",
			emailInput:    "john@example.com",
			passwordInput: "password123",
			setupMock: func(mockRepo *UserRepoMock) {
				// Mock GetByEmail to return ErrNotFound (user doesn't exist)
				mockRepo.On("GetByEmail", mock.Anything, "john@example.com").Return(nil, repo.ErrNotFound)
				// Mock GetByUsername to return an existing user
				existingUser := &entity.User{Username: "johndoe"}
				mockRepo.On("GetByUsername", mock.Anything, "johndoe").Return(existingUser, nil)
			},
			expectedUser:  nil,
			expectedError: repo.ErrDuplicateUsername,
		},
		{
			name:          "DatabaseErrorOnCreate",
			nameInput:     "John Doe",
			usernameInput: "johndoe",
			emailInput:    "john@example.com",
			passwordInput: "password123",
			setupMock: func(mockRepo *UserRepoMock) {
				// Mock GetByEmail to return ErrNotFound (user doesn't exist)
				mockRepo.On("GetByEmail", mock.Anything, "john@example.com").Return(nil, repo.ErrNotFound)
				// Mock GetByUsername to return ErrNotFound (user doesn't exist)
				mockRepo.On("GetByUsername", mock.Anything, "johndoe").Return(nil, repo.ErrNotFound)
				// Mock Create to return an error
				mockRepo.On("Create", mock.Anything, mock.Anything).Return(errors.New("database error"))
			},
			expectedUser:  nil,
			expectedError: errors.New("failed to create user"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock repository
			mockRepo := &UserRepoMock{}

			// Setup mock expectations
			tt.setupMock(mockRepo)

			// Create user service with mock repository
			userService := usecase.NewUserUseCase(mockRepo)

			// Execute the method under test
			user, err := userService.Register(context.Background(), tt.nameInput, tt.usernameInput, tt.emailInput, tt.passwordInput)

			// Assert results
			if tt.expectedError != nil {
				assert.Error(t, err)
				if errors.Is(tt.expectedError, repo.ErrDuplicateEmail) || errors.Is(tt.expectedError, repo.ErrDuplicateUsername) {
					assert.Equal(t, tt.expectedError, err)
				} else {
					assert.Contains(t, err.Error(), tt.expectedError.Error())
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tt.expectedUser.Name, user.Name)
				assert.Equal(t, tt.expectedUser.Username, user.Username)
				assert.Equal(t, tt.expectedUser.Email, user.Email)
				assert.Equal(t, tt.expectedUser.Password, user.Password) // Should be empty
				// Verify that ID was set (not zero value)
				assert.NotEqual(t, uuid.Nil, user.ID)
			}

			// Assert that all expectations were met
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUserService_Login(t *testing.T) {
	// Hash a password for testing
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	// Test cases
	tests := []struct {
		name          string
		emailInput    string
		passwordInput string
		setupMock     func(*UserRepoMock)
		expectedToken string
		expectedError error
	}{
		{
			name:          "Success",
			emailInput:    "john@example.com",
			passwordInput: "password123",
			setupMock: func(mockRepo *UserRepoMock) {
				// Mock GetByEmail to return a user with hashed password
				user := &entity.User{
					ID:       uuid.New(),
					Email:    "john@example.com",
					Password: string(hashedPassword),
				}
				mockRepo.On("GetByEmail", mock.Anything, "john@example.com").Return(user, nil)
			},
			expectedToken: "jwt-token-placeholder",
			expectedError: nil,
		},
		{
			name:          "UserNotFound",
			emailInput:    "john@example.com",
			passwordInput: "password123",
			setupMock: func(mockRepo *UserRepoMock) {
				// Mock GetByEmail to return ErrNotFound
				mockRepo.On("GetByEmail", mock.Anything, "john@example.com").Return(nil, repo.ErrNotFound)
			},
			expectedToken: "",
			expectedError: repo.ErrInvalidCredentials,
		},
		{
			name:          "InvalidPassword",
			emailInput:    "john@example.com",
			passwordInput: "wrongpassword", // Wrong password
			setupMock: func(mockRepo *UserRepoMock) {
				// Mock GetByEmail to return a user with hashed password
				user := &entity.User{
					ID:       uuid.New(),
					Email:    "john@example.com",
					Password: string(hashedPassword),
				}
				mockRepo.On("GetByEmail", mock.Anything, "john@example.com").Return(user, nil)
			},
			expectedToken: "",
			expectedError: repo.ErrInvalidCredentials,
		},
		{
			name:          "DatabaseError",
			emailInput:    "john@example.com",
			passwordInput: "password123",
			setupMock: func(mockRepo *UserRepoMock) {
				// Mock GetByEmail to return an error
				mockRepo.On("GetByEmail", mock.Anything, "john@example.com").Return(nil, errors.New("database error"))
			},
			expectedToken: "",
			expectedError: errors.New("failed to get user"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock repository
			mockRepo := &UserRepoMock{}

			// Setup mock expectations
			tt.setupMock(mockRepo)

			// Create user service with mock repository
			userService := usecase.NewUserUseCase(mockRepo)

			// Execute the method under test
			token, err := userService.Login(context.Background(), tt.emailInput, tt.passwordInput)

			// Assert results
			if tt.expectedError != nil {
				assert.Error(t, err)
				if errors.Is(tt.expectedError, repo.ErrInvalidCredentials) {
					assert.Equal(t, tt.expectedError, err)
				} else {
					assert.Contains(t, err.Error(), tt.expectedError.Error())
				}
			} else {
				assert.NoError(t, err)
				// In a real test, you'd validate the JWT structure, but for this mock, we just check it's not empty
				assert.NotEmpty(t, token)
			}

			// Assert that all expectations were met
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUserService_GetProfile(t *testing.T) {
	// Test cases
	tests := []struct {
		name          string
		usernameInput string
		setupMock     func(*UserRepoMock)
		expectedUser  *entity.User
		expectedError error
	}{
		{
			name:          "Success",
			usernameInput: "johndoe",
			setupMock: func(mockRepo *UserRepoMock) {
				// Mock GetByUsername to return a user
				user := &entity.User{
					ID:       uuid.New(),
					Name:     "John Doe",
					Username: "johndoe",
					Email:    "john@example.com",
					Password: "hashed_password", // Should be cleared in response
				}
				mockRepo.On("GetByUsername", mock.Anything, "johndoe").Return(user, nil)
			},
			expectedUser: &entity.User{
				Name:     "John Doe",
				Username: "johndoe",
				Email:    "john@example.com",
				Password: "", // Should be cleared
			},
			expectedError: nil,
		},
		{
			name:          "UserNotFound",
			usernameInput: "nonexistent",
			setupMock: func(mockRepo *UserRepoMock) {
				// Mock GetByUsername to return ErrNotFound
				mockRepo.On("GetByUsername", mock.Anything, "nonexistent").Return(nil, repo.ErrNotFound)
			},
			expectedUser:  nil,
			expectedError: repo.ErrNotFound,
		},
		{
			name:          "DatabaseError",
			usernameInput: "johndoe",
			setupMock: func(mockRepo *UserRepoMock) {
				// Mock GetByUsername to return an error
				mockRepo.On("GetByUsername", mock.Anything, "johndoe").Return(nil, errors.New("database error"))
			},
			expectedUser:  nil,
			expectedError: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock repository
			mockRepo := &UserRepoMock{}

			// Setup mock expectations
			tt.setupMock(mockRepo)

			// Create user service with mock repository
			userService := usecase.NewUserUseCase(mockRepo)

			// Execute the method under test
			user, err := userService.GetProfile(context.Background(), tt.usernameInput)

			// Assert results
			if tt.expectedError != nil {
				assert.Error(t, err)
				if errors.Is(tt.expectedError, repo.ErrNotFound) {
					assert.Equal(t, tt.expectedError, err)
				} else {
					assert.Contains(t, err.Error(), tt.expectedError.Error())
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tt.expectedUser.Name, user.Name)
				assert.Equal(t, tt.expectedUser.Username, user.Username)
				assert.Equal(t, tt.expectedUser.Email, user.Email)
				assert.Equal(t, tt.expectedUser.Password, user.Password) // Should be empty
				// Verify that ID was set (not zero value)
				assert.NotEqual(t, uuid.Nil, user.ID)
			}

			// Assert that all expectations were met
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUserService_UpdateProfile(t *testing.T) {
	userID := uuid.New()
	name := "John Updated"
	bio := "This is my bio"
	imageURL := "https://example.com/image.jpg"

	// Test cases
	tests := []struct {
		name          string
		userIDInput   uuid.UUID
		nameInput     *string
		bioInput      *string
		imageURLInput *string
		setupMock     func(*UserRepoMock)
		expectedUser  *entity.User
		expectedError error
	}{
		{
			name:          "Success",
			userIDInput:   userID,
			nameInput:     &name,
			bioInput:      &bio,
			imageURLInput: &imageURL,
			setupMock: func(mockRepo *UserRepoMock) {
				// Mock GetByID to return existing user
				existingUser := &entity.User{
					ID:       userID,
					Name:     "John Doe",
					Username: "johndoe",
					Email:    "john@example.com",
					Password: "hashed_password",
				}
				mockRepo.On("GetByID", mock.Anything, userID).Return(existingUser, nil)

				// Mock Update to succeed
				mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(user *entity.User) bool {
					return user.ID == userID && user.Name == name && *user.Bio == bio && *user.ImageURL == imageURL
				})).Return(nil)
			},
			expectedUser: &entity.User{
				ID:       userID,
				Name:     name,
				Username: "johndoe",
				Email:    "john@example.com",
				Password: "", // Should be cleared
			},
			expectedError: nil,
		},
		{
			name:          "UserNotFound",
			userIDInput:   userID,
			nameInput:     &name,
			bioInput:      &bio,
			imageURLInput: &imageURL,
			setupMock: func(mockRepo *UserRepoMock) {
				// Mock GetByID to return ErrNotFound
				mockRepo.On("GetByID", mock.Anything, userID).Return(nil, repo.ErrNotFound)
			},
			expectedUser:  nil,
			expectedError: repo.ErrNotFound,
		},
		{
			name:          "ValidationFailed",
			userIDInput:   userID,
			nameInput:     nil,
			bioInput:      &bio,
			imageURLInput: stringPtr("invalid-url"), // Invalid URL
			setupMock: func(mockRepo *UserRepoMock) {
				// Mock GetByID to return existing user
				existingUser := &entity.User{
					ID:       userID,
					Name:     "John Doe",
					Username: "johndoe",
					Email:    "john@example.com",
					Password: "hashed_password",
				}
				mockRepo.On("GetByID", mock.Anything, userID).Return(existingUser, nil)
			},
			expectedUser:  nil,
			expectedError: errors.New("validation failed"),
		},
		{
			name:          "DatabaseErrorOnUpdate",
			userIDInput:   userID,
			nameInput:     &name,
			bioInput:      &bio,
			imageURLInput: &imageURL,
			setupMock: func(mockRepo *UserRepoMock) {
				// Mock GetByID to return existing user
				existingUser := &entity.User{
					ID:       userID,
					Name:     "John Doe",
					Username: "johndoe",
					Email:    "john@example.com",
					Password: "hashed_password",
				}
				mockRepo.On("GetByID", mock.Anything, userID).Return(existingUser, nil)

				// Mock Update to return an error
				mockRepo.On("Update", mock.Anything, mock.Anything).Return(errors.New("database error"))
			},
			expectedUser:  nil,
			expectedError: errors.New("failed to update profile"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock repository
			mockRepo := &UserRepoMock{}

			// Setup mock expectations
			tt.setupMock(mockRepo)

			// Create user service with mock repository
			userService := usecase.NewUserUseCase(mockRepo)

			// Execute the method under test
			user, err := userService.UpdateProfile(context.Background(), tt.userIDInput, tt.nameInput, tt.bioInput, tt.imageURLInput)

			// Assert results
			if tt.expectedError != nil {
				assert.Error(t, err)
				if errors.Is(tt.expectedError, repo.ErrNotFound) {
					assert.Equal(t, tt.expectedError, err)
				} else {
					assert.Contains(t, err.Error(), tt.expectedError.Error())
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tt.expectedUser.ID, user.ID)
				assert.Equal(t, tt.expectedUser.Name, user.Name)
				assert.Equal(t, tt.expectedUser.Username, user.Username)
				assert.Equal(t, tt.expectedUser.Email, user.Email)
				assert.Equal(t, tt.expectedUser.Password, user.Password) // Should be empty
			}

			// Assert that all expectations were met
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUserService_SearchUsers(t *testing.T) {
	// Test cases
	tests := []struct {
		name          string
		queryInput    string
		limitInput    int
		offsetInput   int
		setupMock     func(*UserRepoMock)
		expectedUsers []entity.User
		expectedError error
	}{
		{
			name:        "Success",
			queryInput:  "john",
			limitInput:  10,
			offsetInput: 0,
			setupMock: func(mockRepo *UserRepoMock) {
				// Mock Search to return users
				users := []entity.User{
					{
						ID:       uuid.New(),
						Name:     "John Doe",
						Username: "johndoe",
						Email:    "john@example.com",
						Password: "hashed_password", // Should be cleared in response
					},
					{
						ID:       uuid.New(),
						Name:     "Johnny Smith",
						Username: "johnnysmith",
						Email:    "johnny@example.com",
						Password: "hashed_password", // Should be cleared in response
					},
				}
				mockRepo.On("Search", mock.Anything, "john", 10, 0).Return(users, nil)
			},
			expectedUsers: []entity.User{
				{
					Name:     "John Doe",
					Username: "johndoe",
					Email:    "john@example.com",
					Password: "", // Should be cleared
				},
				{
					Name:     "Johnny Smith",
					Username: "johnnysmith",
					Email:    "johnny@example.com",
					Password: "", // Should be cleared
				},
			},
			expectedError: nil,
		},
		{
			name:        "NoResults",
			queryInput:  "nonexistent",
			limitInput:  10,
			offsetInput: 0,
			setupMock: func(mockRepo *UserRepoMock) {
				// Mock Search to return empty slice
				mockRepo.On("Search", mock.Anything, "nonexistent", 10, 0).Return([]entity.User{}, nil)
			},
			expectedUsers: []entity.User{},
			expectedError: nil,
		},
		{
			name:        "DatabaseError",
			queryInput:  "john",
			limitInput:  10,
			offsetInput: 0,
			setupMock: func(mockRepo *UserRepoMock) {
				// Mock Search to return an error
				mockRepo.On("Search", mock.Anything, "john", 10, 0).Return(nil, errors.New("database error"))
			},
			expectedUsers: nil,
			expectedError: errors.New("failed to search users"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock repository
			mockRepo := &UserRepoMock{}

			// Setup mock expectations
			tt.setupMock(mockRepo)

			// Create user service with mock repository
			userService := usecase.NewUserUseCase(mockRepo)

			// Execute the method under test
			users, err := userService.SearchUsers(context.Background(), tt.queryInput, tt.limitInput, tt.offsetInput)

			// Assert results
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
				assert.Len(t, users, len(tt.expectedUsers))
				for i, expectedUser := range tt.expectedUsers {
					assert.Equal(t, expectedUser.Name, users[i].Name)
					assert.Equal(t, expectedUser.Username, users[i].Username)
					assert.Equal(t, expectedUser.Email, users[i].Email)
					assert.Equal(t, expectedUser.Password, users[i].Password) // Should be empty
				}
			}

			// Assert that all expectations were met
			mockRepo.AssertExpectations(t)
		})
	}
}

// Helper function to create string pointer
func stringPtr(s string) *string {
	return &s
}
