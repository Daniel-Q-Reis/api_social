package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"social/api/internal/entity"
	"social/api/internal/repo"
	"social/api/internal/usecase"
)

func TestInteractionService_LikePost(t *testing.T) {
	postID := uuid.New()
	userID := uuid.New()

	// Test cases
	tests := []struct {
		name          string
		postID        uuid.UUID
		userID        uuid.UUID
		setupMock     func(*LikeRepoMock, *FollowRepoMock, *UserRepoMock)
		expectedError error
	}{
		{
			name:   "Success",
			postID: postID,
			userID: userID,
			setupMock: func(mockLikeRepo *LikeRepoMock, _, _ *UserRepoMock) {
				// Mock Create to succeed
				mockLikeRepo.On("Create", mock.Anything, mock.MatchedBy(func(like *entity.Like) bool {
					return like.PostID == postID && like.UserID == userID
				})).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:   "DatabaseError",
			postID: postID,
			userID: userID,
			setupMock: func(mockLikeRepo *LikeRepoMock, _, _ *UserRepoMock) {
				// Mock Create to return an error
				mockLikeRepo.On("Create", mock.Anything, mock.Anything).Return(errors.New("database error"))
			},
			expectedError: errors.New("failed to like post"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock repositories
			mockLikeRepo := &LikeRepoMock{}
			mockFollowRepo := &FollowRepoMock{}
			mockUserRepo := &UserRepoMock{}

			// Setup mock expectations
			tt.setupMock(mockLikeRepo, mockFollowRepo, mockUserRepo)

			// Create interaction service with mock repositories
			interactionService := usecase.NewInteractionUseCase(mockLikeRepo, mockFollowRepo, mockUserRepo)

			// Execute the method under test
			err := interactionService.LikePost(context.Background(), tt.postID, tt.userID)

			// Assert results
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}

			// Assert that all expectations were met
			mockLikeRepo.AssertExpectations(t)
			mockFollowRepo.AssertExpectations(t)
			mockUserRepo.AssertExpectations(t)
		})
	}
}

func TestInteractionService_UnlikePost(t *testing.T) {
	postID := uuid.New()
	userID := uuid.New()

	// Test cases
	tests := []struct {
		name          string
		postID        uuid.UUID
		userID        uuid.UUID
		setupMock     func(*LikeRepoMock, *FollowRepoMock, *UserRepoMock)
		expectedError error
	}{
		{
			name:   "Success",
			postID: postID,
			userID: userID,
			setupMock: func(mockLikeRepo *LikeRepoMock, _, _ *UserRepoMock) {
				// Mock Delete to succeed
				mockLikeRepo.On("Delete", mock.Anything, userID, postID).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:   "DatabaseError",
			postID: postID,
			userID: userID,
			setupMock: func(mockLikeRepo *LikeRepoMock, _, _ *UserRepoMock) {
				// Mock Delete to return an error
				mockLikeRepo.On("Delete", mock.Anything, userID, postID).Return(errors.New("database error"))
			},
			expectedError: errors.New("failed to unlike post"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock repositories
			mockLikeRepo := &LikeRepoMock{}
			mockFollowRepo := &FollowRepoMock{}
			mockUserRepo := &UserRepoMock{}

			// Setup mock expectations
			tt.setupMock(mockLikeRepo, mockFollowRepo, mockUserRepo)

			// Create interaction service with mock repositories
			interactionService := usecase.NewInteractionUseCase(mockLikeRepo, mockFollowRepo, mockUserRepo)

			// Execute the method under test
			err := interactionService.UnlikePost(context.Background(), tt.postID, tt.userID)

			// Assert results
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}

			// Assert that all expectations were met
			mockLikeRepo.AssertExpectations(t)
			mockFollowRepo.AssertExpectations(t)
			mockUserRepo.AssertExpectations(t)
		})
	}
}

func TestInteractionService_FollowUser(t *testing.T) {
	userID := uuid.New()
	followerID := uuid.New()

	// Test cases
	tests := []struct {
		name          string
		userID        uuid.UUID
		followerID    uuid.UUID
		setupMock     func(*LikeRepoMock, *FollowRepoMock, *UserRepoMock)
		expectedError error
	}{
		{
			name:       "Success",
			userID:     userID,
			followerID: followerID,
			setupMock: func(_ *LikeRepoMock, mockFollowRepo *FollowRepoMock, _ *UserRepoMock) {
				// Mock Create to succeed
				mockFollowRepo.On("Create", mock.Anything, mock.MatchedBy(func(follow *entity.Follow) bool {
					return follow.UserID == userID && follow.FollowerID == followerID
				})).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:       "CannotFollowSelf",
			userID:     userID,
			followerID: userID, // Same as userID
			setupMock: func(_, _, _ *UserRepoMock) {
				// No need to mock repository as validation happens before repository calls
			},
			expectedError: errors.New("you cannot follow yourself"),
		},
		{
			name:       "DatabaseError",
			userID:     userID,
			followerID: followerID,
			setupMock: func(_ *LikeRepoMock, mockFollowRepo *FollowRepoMock, _ *UserRepoMock) {
				// Mock Create to return an error
				mockFollowRepo.On("Create", mock.Anything, mock.Anything).Return(errors.New("database error"))
			},
			expectedError: errors.New("failed to follow user"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock repositories
			mockLikeRepo := &LikeRepoMock{}
			mockFollowRepo := &FollowRepoMock{}
			mockUserRepo := &UserRepoMock{}

			// Setup mock expectations
			tt.setupMock(mockLikeRepo, mockFollowRepo, mockUserRepo)

			// Create interaction service with mock repositories
			interactionService := usecase.NewInteractionUseCase(mockLikeRepo, mockFollowRepo, mockUserRepo)

			// Execute the method under test
			err := interactionService.FollowUser(context.Background(), tt.userID, tt.followerID)

			// Assert results
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}

			// Assert that all expectations were met
			mockLikeRepo.AssertExpectations(t)
			mockFollowRepo.AssertExpectations(t)
			mockUserRepo.AssertExpectations(t)
		})
	}
}

func TestInteractionService_UnfollowUser(t *testing.T) {
	userID := uuid.New()
	followerID := uuid.New()

	// Test cases
	tests := []struct {
		name          string
		userID        uuid.UUID
		followerID    uuid.UUID
		setupMock     func(*LikeRepoMock, *FollowRepoMock, *UserRepoMock)
		expectedError error
	}{
		{
			name:       "Success",
			userID:     userID,
			followerID: followerID,
			setupMock: func(_, mockFollowRepo *FollowRepoMock, _) {
				// Mock Delete to succeed
				mockFollowRepo.On("Delete", mock.Anything, userID, followerID).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:       "DatabaseError",
			userID:     userID,
			followerID: followerID,
			setupMock: func(_, mockFollowRepo *FollowRepoMock, _) {
				// Mock Delete to return an error
				mockFollowRepo.On("Delete", mock.Anything, userID, followerID).Return(errors.New("database error"))
			},
			expectedError: errors.New("failed to unfollow user"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock repositories
			mockLikeRepo := &LikeRepoMock{}
			mockFollowRepo := &FollowRepoMock{}
			mockUserRepo := &UserRepoMock{}

			// Setup mock expectations
			tt.setupMock(mockLikeRepo, mockFollowRepo, mockUserRepo)

			// Create interaction service with mock repositories
			interactionService := usecase.NewInteractionUseCase(mockLikeRepo, mockFollowRepo, mockUserRepo)

			// Execute the method under test
			err := interactionService.UnfollowUser(context.Background(), tt.userID, tt.followerID)

			// Assert results
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}

			// Assert that all expectations were met
			mockLikeRepo.AssertExpectations(t)
			mockFollowRepo.AssertExpectations(t)
			mockUserRepo.AssertExpectations(t)
		})
	}
}

func TestInteractionService_GetFollowers(t *testing.T) {
	username := "johndoe"
	userID := uuid.New()
	followerID1 := uuid.New()
	followerID2 := uuid.New()

	// Test cases
	tests := []struct {
		name             string
		username         string
		limit            int
		offset           int
		setupMock        func(*LikeRepoMock, *FollowRepoMock, *UserRepoMock)
		expectedFollowers []entity.User
		expectedError     error
	}{
		{
			name:     "Success",
			username: username,
			limit:    10,
			offset:   0,
			setupMock: func(_, _ *FollowRepoMock, mockUserRepo *UserRepoMock) {
				// Mock GetByUsername to return a user
				user := &entity.User{
					ID:       userID,
					Username: username,
				}
				mockUserRepo.On("GetByUsername", mock.Anything, username).Return(user, nil)

				// Mock GetFollowers to return followers
				followers := []entity.User{
					{
						ID:       followerID1,
						Name:     "Follower One",
						Username: "follower1",
						Email:    "follower1@example.com",
						Password: "hashed_password", // Should be cleared in response
					},
					{
						ID:       followerID2,
						Name:     "Follower Two",
						Username: "follower2",
						Email:    "follower2@example.com",
						Password: "hashed_password", // Should be cleared in response
					},
				}
				mockFollowRepo.On("GetFollowers", mock.Anything, userID, 10, 0).Return(followers, nil)
			},
			expectedFollowers: []entity.User{
				{
					ID:       followerID1,
					Name:     "Follower One",
					Username: "follower1",
					Email:    "follower1@example.com",
					Password: "", // Should be cleared
				},
				{
					ID:       followerID2,
					Name:     "Follower Two",
					Username: "follower2",
					Email:    "follower2@example.com",
					Password: "", // Should be cleared
				},
			},
			expectedError: nil,
		},
		{
			name:     "UserNotFound",
			username: "nonexistent",
			limit:    10,
			offset:   0,
			setupMock: func(_, _, mockUserRepo *UserRepoMock) {
				// Mock GetByUsername to return ErrNotFound
				mockUserRepo.On("GetByUsername", mock.Anything, "nonexistent").Return(nil, repo.ErrNotFound)
			},
			expectedFollowers: nil,
			expectedError:     repo.ErrNotFound,
		},
		{
			name:     "DatabaseErrorOnGetUser",
			username: username,
			limit:    10,
			offset:   0,
			setupMock: func(_, _, mockUserRepo *UserRepoMock) {
				// Mock GetByUsername to return an error
				mockUserRepo.On("GetByUsername", mock.Anything, username).Return(nil, errors.New("database error"))
			},
			expectedFollowers: nil,
			expectedError:     errors.New("database error"),
		},
		{
			name:     "DatabaseErrorOnGetFollowers",
			username: username,
			limit:    10,
			offset:   0,
			setupMock: func(mockLikeRepo *LikeRepoMock, mockFollowRepo *FollowRepoMock, mockUserRepo *UserRepoMock) {
				// Mock GetByUsername to return a user
				user := &entity.User{
					ID:       userID,
					Username: username,
				}
				mockUserRepo.On("GetByUsername", mock.Anything, username).Return(user, nil)

				// Mock GetFollowers to return an error
				mockFollowRepo.On("GetFollowers", mock.Anything, userID, 10, 0).Return(nil, errors.New("database error"))
			},
			expectedFollowers: nil,
			expectedError:     errors.New("failed to get followers"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock repositories
			mockLikeRepo := &LikeRepoMock{}
			mockFollowRepo := &FollowRepoMock{}
			mockUserRepo := &UserRepoMock{}

			// Setup mock expectations
			tt.setupMock(mockLikeRepo, mockFollowRepo, mockUserRepo)

			// Create interaction service with mock repositories
			interactionService := usecase.NewInteractionUseCase(mockLikeRepo, mockFollowRepo, mockUserRepo)

			// Execute the method under test
			followers, err := interactionService.GetFollowers(context.Background(), tt.username, tt.limit, tt.offset)

			// Assert results
			if tt.expectedError != nil {
				assert.Error(t, err)
				if tt.expectedError == repo.ErrNotFound {
					assert.Equal(t, tt.expectedError, err)
				} else {
					assert.Contains(t, err.Error(), tt.expectedError.Error())
				}
			} else {
				assert.NoError(t, err)
				assert.Len(t, followers, len(tt.expectedFollowers))
				for i, expectedFollower := range tt.expectedFollowers {
					assert.Equal(t, expectedFollower.ID, followers[i].ID)
					assert.Equal(t, expectedFollower.Name, followers[i].Name)
					assert.Equal(t, expectedFollower.Username, followers[i].Username)
					assert.Equal(t, expectedFollower.Email, followers[i].Email)
					assert.Equal(t, expectedFollower.Password, followers[i].Password) // Should be empty
				}
			}

			// Assert that all expectations were met
			mockLikeRepo.AssertExpectations(t)
			mockFollowRepo.AssertExpectations(t)
			mockUserRepo.AssertExpectations(t)
		})
	}
}

func TestInteractionService_GetFollowing(t *testing.T) {
	username := "johndoe"
	userID := uuid.New()
	followingID1 := uuid.New()
	followingID2 := uuid.New()

	// Test cases
	tests := []struct {
		name              string
		username          string
		limit             int
		offset            int
		setupMock         func(*LikeRepoMock, *FollowRepoMock, *UserRepoMock)
		expectedFollowing []entity.User
		expectedError     error
	}{
		{
			name:     "Success",
			username: username,
			limit:    10,
			offset:   0,
			setupMock: func(_, _ *FollowRepoMock, mockUserRepo *UserRepoMock) {
				// Mock GetByUsername to return a user
				user := &entity.User{
					ID:       userID,
					Username: username,
				}
				mockUserRepo.On("GetByUsername", mock.Anything, username).Return(user, nil)

				// Mock GetFollowing to return following
				following := []entity.User{
					{
						ID:       followingID1,
						Name:     "Following One",
						Username: "following1",
						Email:    "following1@example.com",
						Password: "hashed_password", // Should be cleared in response
					},
					{
						ID:       followingID2,
						Name:     "Following Two",
						Username: "following2",
						Email:    "following2@example.com",
						Password: "hashed_password", // Should be cleared in response
					},
				}
				mockFollowRepo.On("GetFollowing", mock.Anything, userID, 10, 0).Return(following, nil)
			},
			expectedFollowing: []entity.User{
				{
					ID:       followingID1,
					Name:     "Following One",
					Username: "following1",
					Email:    "following1@example.com",
					Password: "", // Should be cleared
				},
				{
					ID:       followingID2,
					Name:     "Following Two",
					Username: "following2",
					Email:    "following2@example.com",
					Password: "", // Should be cleared
				},
			},
			expectedError: nil,
		},
		{
			name:     "UserNotFound",
			username: "nonexistent",
			limit:    10,
			offset:   0,
			setupMock: func(_, _, mockUserRepo *UserRepoMock) {
				// Mock GetByUsername to return ErrNotFound
				mockUserRepo.On("GetByUsername", mock.Anything, "nonexistent").Return(nil, repo.ErrNotFound)
			},
			expectedFollowing: nil,
			expectedError:     repo.ErrNotFound,
		},
		{
			name:     "DatabaseErrorOnGetUser",
			username: username,
			limit:    10,
			offset:   0,
			setupMock: func(_, _, mockUserRepo *UserRepoMock) {
				// Mock GetByUsername to return an error
				mockUserRepo.On("GetByUsername", mock.Anything, username).Return(nil, errors.New("database error"))
			},
			expectedFollowing: nil,
			expectedError:     errors.New("database error"),
		},
		{
			name:     "DatabaseErrorOnGetFollowing",
			username: username,
			limit:    10,
			offset:   0,
			setupMock: func(mockLikeRepo *LikeRepoMock, mockFollowRepo *FollowRepoMock, mockUserRepo *UserRepoMock) {
				// Mock GetByUsername to return a user
				user := &entity.User{
					ID:       userID,
					Username: username,
				}
				mockUserRepo.On("GetByUsername", mock.Anything, username).Return(user, nil)

				// Mock GetFollowing to return an error
				mockFollowRepo.On("GetFollowing", mock.Anything, userID, 10, 0).Return(nil, errors.New("database error"))
			},
			expectedFollowing: nil,
			expectedError:     errors.New("failed to get following"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock repositories
			mockLikeRepo := &LikeRepoMock{}
			mockFollowRepo := &FollowRepoMock{}
			mockUserRepo := &UserRepoMock{}

			// Setup mock expectations
			tt.setupMock(mockLikeRepo, mockFollowRepo, mockUserRepo)

			// Create interaction service with mock repositories
			interactionService := usecase.NewInteractionUseCase(mockLikeRepo, mockFollowRepo, mockUserRepo)

			// Execute the method under test
			following, err := interactionService.GetFollowing(context.Background(), tt.username, tt.limit, tt.offset)

			// Assert results
			if tt.expectedError != nil {
				assert.Error(t, err)
				if tt.expectedError == repo.ErrNotFound {
					assert.Equal(t, tt.expectedError, err)
				} else {
					assert.Contains(t, err.Error(), tt.expectedError.Error())
				}
			} else {
				assert.NoError(t, err)
				assert.Len(t, following, len(tt.expectedFollowing))
				for i, expectedFollowing := range tt.expectedFollowing {
					assert.Equal(t, expectedFollowing.ID, following[i].ID)
					assert.Equal(t, expectedFollowing.Name, following[i].Name)
					assert.Equal(t, expectedFollowing.Username, following[i].Username)
					assert.Equal(t, expectedFollowing.Email, following[i].Email)
					assert.Equal(t, expectedFollowing.Password, following[i].Password) // Should be empty
				}
			}

			// Assert that all expectations were met
			mockLikeRepo.AssertExpectations(t)
			mockFollowRepo.AssertExpectations(t)
			mockUserRepo.AssertExpectations(t)
		})
	}
}