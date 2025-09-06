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
)

func TestPostService_CreatePost(t *testing.T) {
	authorID := uuid.New()
	content := "This is a test post"
	imageURL := "https://example.com/image.jpg"

	// Test cases
	tests := []struct {
		name          string
		authorID      uuid.UUID
		content       string
		imageURL      *string
		setupMock     func(*PostRepoMock, *UserRepoMock)
		expectedPost  *entity.Post
		expectedError error
	}{
		{
			name:     "Success",
			authorID: authorID,
			content:  content,
			imageURL: &imageURL,
			setupMock: func(mockPostRepo *PostRepoMock, _ *UserRepoMock) {
				// Mock Create to succeed
				mockPostRepo.On("Create", mock.Anything, mock.MatchedBy(func(post *entity.Post) bool {
					return post.AuthorID == authorID && post.Content == content && *post.ImageURL == imageURL
				})).Return(nil)
			},
			expectedPost: &entity.Post{
				AuthorID: authorID,
				Content:  content,
				ImageURL: &imageURL,
			},
			expectedError: nil,
		},
		{
			name:     "ValidationFailed_EmptyContent",
			authorID: authorID,
			content:  "", // Empty content
			imageURL: &imageURL,
			setupMock: func(_ *PostRepoMock, _ *UserRepoMock) {
				// No need to mock repository as validation happens before repository calls
			},
			expectedPost:  nil,
			expectedError: errors.New("validation failed"),
		},
		{
			name:     "ValidationFailed_InvalidImageURL",
			authorID: authorID,
			content:  content,
			imageURL: stringPtr("invalid-url"), // Invalid URL
			setupMock: func(_ *PostRepoMock, _ *UserRepoMock) {
				// No need to mock repository as validation happens before repository calls
			},
			expectedPost:  nil,
			expectedError: errors.New("validation failed"),
		},
		{
			name:     "DatabaseError",
			authorID: authorID,
			content:  content,
			imageURL: &imageURL,
			setupMock: func(mockPostRepo *PostRepoMock, _ *UserRepoMock) {
				// Mock Create to return an error
				mockPostRepo.On("Create", mock.Anything, mock.Anything).Return(errors.New("database error"))
			},
			expectedPost:  nil,
			expectedError: errors.New("failed to create post"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock repositories
			mockPostRepo := &PostRepoMock{}
			mockUserRepo := &UserRepoMock{}

			// Setup mock expectations
			tt.setupMock(mockPostRepo, mockUserRepo)

			// Create post service with mock repositories
			postService := usecase.NewPostUseCase(mockPostRepo, mockUserRepo)

			// Execute the method under test
			post, err := postService.CreatePost(context.Background(), tt.authorID, tt.content, tt.imageURL)

			// Assert results
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, post)
				assert.Equal(t, tt.expectedPost.AuthorID, post.AuthorID)
				assert.Equal(t, tt.expectedPost.Content, post.Content)
				assert.Equal(t, *tt.expectedPost.ImageURL, *post.ImageURL)
				// Verify that ID was set (not zero value)
				assert.NotEqual(t, uuid.Nil, post.ID)
			}

			// Assert that all expectations were met
			mockPostRepo.AssertExpectations(t)
			mockUserRepo.AssertExpectations(t)
		})
	}
}

func TestPostService_GetPostByID(t *testing.T) {
	postID := uuid.New()
	authorID := uuid.New()

	// Test cases
	tests := []struct {
		name          string
		postID        uuid.UUID
		setupMock     func(*PostRepoMock, *UserRepoMock)
		expectedPost  *entity.Post
		expectedError error
	}{
		{
			name:   "Success",
			postID: postID,
			setupMock: func(mockPostRepo *PostRepoMock, _ *UserRepoMock) {
				// Mock GetByID to return a post
				post := &entity.Post{
					ID:       postID,
					AuthorID: authorID,
					Content:  "This is a test post",
				}
				mockPostRepo.On("GetByID", mock.Anything, postID).Return(post, nil)
			},
			expectedPost: &entity.Post{
				ID:       postID,
				AuthorID: authorID,
				Content:  "This is a test post",
			},
			expectedError: nil,
		},
		{
			name:   "PostNotFound",
			postID: postID,
			setupMock: func(mockPostRepo *PostRepoMock, _ *UserRepoMock) {
				// Mock GetByID to return ErrNotFound
				mockPostRepo.On("GetByID", mock.Anything, postID).Return(nil, repo.ErrNotFound)
			},
			expectedPost:  nil,
			expectedError: repo.ErrNotFound,
		},
		{
			name:   "DatabaseError",
			postID: postID,
			setupMock: func(mockPostRepo *PostRepoMock, _ *UserRepoMock) {
				// Mock GetByID to return an error
				mockPostRepo.On("GetByID", mock.Anything, postID).Return(nil, errors.New("database error"))
			},
			expectedPost:  nil,
			expectedError: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock repositories
			mockPostRepo := &PostRepoMock{}
			mockUserRepo := &UserRepoMock{}

			// Setup mock expectations
			tt.setupMock(mockPostRepo, mockUserRepo)

			// Create post service with mock repositories
			postService := usecase.NewPostUseCase(mockPostRepo, mockUserRepo)

			// Execute the method under test
			post, err := postService.GetPostByID(context.Background(), tt.postID)

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
				assert.NotNil(t, post)
				assert.Equal(t, tt.expectedPost.ID, post.ID)
				assert.Equal(t, tt.expectedPost.AuthorID, post.AuthorID)
				assert.Equal(t, tt.expectedPost.Content, post.Content)
			}

			// Assert that all expectations were met
			mockPostRepo.AssertExpectations(t)
			mockUserRepo.AssertExpectations(t)
		})
	}
}

func TestPostService_GetPostsByUser(t *testing.T) {
	username := "johndoe"
	authorID := uuid.New()
	postID1 := uuid.New()
	postID2 := uuid.New()

	// Test cases
	tests := []struct {
		name          string
		username      string
		limit         int
		offset        int
		setupMock     func(*PostRepoMock, *UserRepoMock)
		expectedPosts []entity.Post
		expectedError error
	}{
		{
			name:     "Success",
			username: username,
			limit:    10,
			offset:   0,
			setupMock: func(mockPostRepo *PostRepoMock, mockUserRepo *UserRepoMock) {
				// Mock GetByUsername to return a user
				user := &entity.User{
					ID:       authorID,
					Username: username,
				}
				mockUserRepo.On("GetByUsername", mock.Anything, username).Return(user, nil)

				// Mock GetByAuthorID to return posts
				posts := []entity.Post{
					{
						ID:       postID1,
						AuthorID: authorID,
						Content:  "First post",
					},
					{
						ID:       postID2,
						AuthorID: authorID,
						Content:  "Second post",
					},
				}
				mockPostRepo.On("GetByAuthorID", mock.Anything, authorID, 10, 0).Return(posts, nil)
			},
			expectedPosts: []entity.Post{
				{
					ID:       postID1,
					AuthorID: authorID,
					Content:  "First post",
				},
				{
					ID:       postID2,
					AuthorID: authorID,
					Content:  "Second post",
				},
			},
			expectedError: nil,
		},
		{
			name:     "UserNotFound",
			username: "nonexistent",
			limit:    10,
			offset:   0,
			setupMock: func(_ *PostRepoMock, mockUserRepo *UserRepoMock) {
				// Mock GetByUsername to return ErrNotFound
				mockUserRepo.On("GetByUsername", mock.Anything, "nonexistent").Return(nil, repo.ErrNotFound)
			},
			expectedPosts: nil,
			expectedError: repo.ErrNotFound,
		},
		{
			name:     "DatabaseErrorOnGetUser",
			username: username,
			limit:    10,
			offset:   0,
			setupMock: func(_ *PostRepoMock, mockUserRepo *UserRepoMock) {
				// Mock GetByUsername to return an error
				mockUserRepo.On("GetByUsername", mock.Anything, username).Return(nil, errors.New("database error"))
			},
			expectedPosts: nil,
			expectedError: errors.New("database error"),
		},
		{
			name:     "DatabaseErrorOnGetPosts",
			username: username,
			limit:    10,
			offset:   0,
			setupMock: func(mockPostRepo *PostRepoMock, mockUserRepo *UserRepoMock) {
				// Mock GetByUsername to return a user
				user := &entity.User{
					ID:       authorID,
					Username: username,
				}
				mockUserRepo.On("GetByUsername", mock.Anything, username).Return(user, nil)

				// Mock GetByAuthorID to return an error
				mockPostRepo.On("GetByAuthorID", mock.Anything, authorID, 10, 0).Return(nil, errors.New("database error"))
			},
			expectedPosts: nil,
			expectedError: errors.New("failed to get posts"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock repositories
			mockPostRepo := &PostRepoMock{}
			mockUserRepo := &UserRepoMock{}

			// Setup mock expectations
			tt.setupMock(mockPostRepo, mockUserRepo)

			// Create post service with mock repositories
			postService := usecase.NewPostUseCase(mockPostRepo, mockUserRepo)

			// Execute the method under test
			posts, err := postService.GetPostsByUser(context.Background(), tt.username, tt.limit, tt.offset)

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
				assert.Len(t, posts, len(tt.expectedPosts))
				for i, expectedPost := range tt.expectedPosts {
					assert.Equal(t, expectedPost.ID, posts[i].ID)
					assert.Equal(t, expectedPost.AuthorID, posts[i].AuthorID)
					assert.Equal(t, expectedPost.Content, posts[i].Content)
				}
			}

			// Assert that all expectations were met
			mockPostRepo.AssertExpectations(t)
			mockUserRepo.AssertExpectations(t)
		})
	}
}

func TestPostService_UpdatePost(t *testing.T) {
	postID := uuid.New()
	userID := uuid.New()
	authorID := userID // Same user (authorized)
	content := "Updated content"
	imageURL := "https://example.com/updated-image.jpg"

	// Test cases
	tests := []struct {
		name          string
		postID        uuid.UUID
		userID        uuid.UUID
		content       string
		imageURL      *string
		setupMock     func(*PostRepoMock, *UserRepoMock)
		expectedPost  *entity.Post
		expectedError error
	}{
		{
			name:     "Success",
			postID:   postID,
			userID:   userID,
			content:  content,
			imageURL: &imageURL,
			setupMock: func(mockPostRepo *PostRepoMock, _ *UserRepoMock) {
				// Mock GetByID to return a post
				post := &entity.Post{
					ID:       postID,
					AuthorID: authorID, // Same as userID
					Content:  "Original content",
				}
				mockPostRepo.On("GetByID", mock.Anything, postID).Return(post, nil)

				// Mock Update to succeed
				mockPostRepo.On("Update", mock.Anything, mock.MatchedBy(func(post *entity.Post) bool {
					return post.ID == postID && post.Content == content && *post.ImageURL == imageURL
				})).Return(nil)
			},
			expectedPost: &entity.Post{
				ID:       postID,
				AuthorID: authorID,
				Content:  content,
				ImageURL: &imageURL,
			},
			expectedError: nil,
		},
		{
			name:     "PostNotFound",
			postID:   postID,
			userID:   userID,
			content:  content,
			imageURL: &imageURL,
			setupMock: func(mockPostRepo *PostRepoMock, _ *UserRepoMock) {
				// Mock GetByID to return ErrNotFound
				mockPostRepo.On("GetByID", mock.Anything, postID).Return(nil, repo.ErrNotFound)
			},
			expectedPost:  nil,
			expectedError: repo.ErrNotFound,
		},
		{
			name:     "Unauthorized",
			postID:   postID,
			userID:   userID,
			content:  content,
			imageURL: &imageURL,
			setupMock: func(mockPostRepo *PostRepoMock, _ *UserRepoMock) {
				// Mock GetByID to return a post with different author
				post := &entity.Post{
					ID:       postID,
					AuthorID: uuid.New(), // Different author
					Content:  "Original content",
				}
				mockPostRepo.On("GetByID", mock.Anything, postID).Return(post, nil)
			},
			expectedPost:  nil,
			expectedError: repo.ErrUnauthorized,
		},
		{
			name:     "ValidationFailed_EmptyContent",
			postID:   postID,
			userID:   userID,
			content:  "", // Empty content
			imageURL: &imageURL,
			setupMock: func(mockPostRepo *PostRepoMock, _ *UserRepoMock) {
				// Mock GetByID to return a post
				post := &entity.Post{
					ID:       postID,
					AuthorID: authorID, // Same as userID
					Content:  "Original content",
				}
				mockPostRepo.On("GetByID", mock.Anything, postID).Return(post, nil)
			},
			expectedPost:  nil,
			expectedError: errors.New("validation failed"),
		},
		{
			name:     "DatabaseErrorOnUpdate",
			postID:   postID,
			userID:   userID,
			content:  content,
			imageURL: &imageURL,
			setupMock: func(mockPostRepo *PostRepoMock, _ *UserRepoMock) {
				// Mock GetByID to return a post
				post := &entity.Post{
					ID:       postID,
					AuthorID: authorID, // Same as userID
					Content:  "Original content",
				}
				mockPostRepo.On("GetByID", mock.Anything, postID).Return(post, nil)

				// Mock Update to return an error
				mockPostRepo.On("Update", mock.Anything, mock.Anything).Return(errors.New("database error"))
			},
			expectedPost:  nil,
			expectedError: errors.New("failed to update post"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock repositories
			mockPostRepo := &PostRepoMock{}
			mockUserRepo := &UserRepoMock{}

			// Setup mock expectations
			tt.setupMock(mockPostRepo, mockUserRepo)

			// Create post service with mock repositories
			postService := usecase.NewPostUseCase(mockPostRepo, mockUserRepo)

			// Execute the method under test
			post, err := postService.UpdatePost(context.Background(), tt.postID, tt.userID, tt.content, tt.imageURL)

			// Assert results
			if tt.expectedError != nil {
				assert.Error(t, err)
				if errors.Is(tt.expectedError, repo.ErrNotFound) || errors.Is(tt.expectedError, repo.ErrUnauthorized) {
					assert.Equal(t, tt.expectedError, err)
				} else {
					assert.Contains(t, err.Error(), tt.expectedError.Error())
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, post)
				assert.Equal(t, tt.expectedPost.ID, post.ID)
				assert.Equal(t, tt.expectedPost.AuthorID, post.AuthorID)
				assert.Equal(t, tt.expectedPost.Content, post.Content)
				assert.Equal(t, *tt.expectedPost.ImageURL, *post.ImageURL)
			}

			// Assert that all expectations were met
			mockPostRepo.AssertExpectations(t)
			mockUserRepo.AssertExpectations(t)
		})
	}
}

func TestPostService_DeletePost(t *testing.T) {
	postID := uuid.New()
	userID := uuid.New()
	authorID := userID // Same user (authorized)

	// Test cases
	tests := []struct {
		name          string
		postID        uuid.UUID
		userID        uuid.UUID
		setupMock     func(*PostRepoMock, *UserRepoMock)
		expectedError error
	}{
		{
			name:   "Success",
			postID: postID,
			userID: userID,
			setupMock: func(mockPostRepo *PostRepoMock, _ *UserRepoMock) {
				// Mock GetByID to return a post
				post := &entity.Post{
					ID:       postID,
					AuthorID: authorID, // Same as userID
				}
				mockPostRepo.On("GetByID", mock.Anything, postID).Return(post, nil)

				// Mock Delete to succeed
				mockPostRepo.On("Delete", mock.Anything, postID).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:   "PostNotFound",
			postID: postID,
			userID: userID,
			setupMock: func(mockPostRepo *PostRepoMock, _ *UserRepoMock) {
				// Mock GetByID to return ErrNotFound
				mockPostRepo.On("GetByID", mock.Anything, postID).Return(nil, repo.ErrNotFound)
			},
			expectedError: repo.ErrNotFound,
		},
		{
			name:   "Unauthorized",
			postID: postID,
			userID: userID,
			setupMock: func(mockPostRepo *PostRepoMock, _ *UserRepoMock) {
				// Mock GetByID to return a post with different author
				post := &entity.Post{
					ID:       postID,
					AuthorID: uuid.New(), // Different author
				}
				mockPostRepo.On("GetByID", mock.Anything, postID).Return(post, nil)
			},
			expectedError: repo.ErrUnauthorized,
		},
		{
			name:   "DatabaseError",
			postID: postID,
			userID: userID,
			setupMock: func(mockPostRepo *PostRepoMock, _ *UserRepoMock) {
				// Mock GetByID to return a post
				post := &entity.Post{
					ID:       postID,
					AuthorID: authorID, // Same as userID
				}
				mockPostRepo.On("GetByID", mock.Anything, postID).Return(post, nil)

				// Mock Delete to return an error
				mockPostRepo.On("Delete", mock.Anything, postID).Return(errors.New("database error"))
			},
			expectedError: errors.New("failed to delete post"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock repositories
			mockPostRepo := &PostRepoMock{}
			mockUserRepo := &UserRepoMock{}

			// Setup mock expectations
			tt.setupMock(mockPostRepo, mockUserRepo)

			// Create post service with mock repositories
			postService := usecase.NewPostUseCase(mockPostRepo, mockUserRepo)

			// Execute the method under test
			err := postService.DeletePost(context.Background(), tt.postID, tt.userID)

			// Assert results
			if tt.expectedError != nil {
				assert.Error(t, err)
				if errors.Is(tt.expectedError, repo.ErrNotFound) || errors.Is(tt.expectedError, repo.ErrUnauthorized) {
					assert.Equal(t, tt.expectedError, err)
				} else {
					assert.Contains(t, err.Error(), tt.expectedError.Error())
				}
			} else {
				assert.NoError(t, err)
			}

			// Assert that all expectations were met
			mockPostRepo.AssertExpectations(t)
			mockUserRepo.AssertExpectations(t)
		})
	}
}

func TestPostService_GetFeed(t *testing.T) {
	userID := uuid.New()

	// Test cases
	tests := []struct {
		name          string
		userID        uuid.UUID
		limit         int
		offset        int
		setupMock     func(*PostRepoMock, *UserRepoMock)
		expectedPosts []entity.Post
		expectedError error
	}{
		{
			name:   "Success",
			userID: userID,
			limit:  10,
			offset: 0,
			setupMock: func(mockPostRepo *PostRepoMock, _ *UserRepoMock) {
				// Mock GetFeed to return posts
				posts := []entity.Post{
					{
						Content: "First post in feed",
					},
					{
						Content: "Second post in feed",
					},
				}
				mockPostRepo.On("GetFeed", mock.Anything, userID, 10, 0).Return(posts, nil)
			},
			expectedPosts: []entity.Post{
				{
					Content: "First post in feed",
				},
				{
					Content: "Second post in feed",
				},
			},
			expectedError: nil,
		},
		{
			name:   "NoPostsInFeed",
			userID: userID,
			limit:  10,
			offset: 0,
			setupMock: func(mockPostRepo *PostRepoMock, _ *UserRepoMock) {
				// Mock GetFeed to return empty slice
				mockPostRepo.On("GetFeed", mock.Anything, userID, 10, 0).Return([]entity.Post{}, nil)
			},
			expectedPosts: []entity.Post{},
			expectedError: nil,
		},
		{
			name:   "DatabaseError",
			userID: userID,
			limit:  10,
			offset: 0,
			setupMock: func(mockPostRepo *PostRepoMock, _ *UserRepoMock) {
				// Mock GetFeed to return an error
				mockPostRepo.On("GetFeed", mock.Anything, userID, 10, 0).Return(nil, errors.New("database error"))
			},
			expectedPosts: nil,
			expectedError: errors.New("failed to get feed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock repositories
			mockPostRepo := &PostRepoMock{}
			mockUserRepo := &UserRepoMock{}

			// Setup mock expectations
			tt.setupMock(mockPostRepo, mockUserRepo)

			// Create post service with mock repositories
			postService := usecase.NewPostUseCase(mockPostRepo, mockUserRepo)

			// Execute the method under test
			posts, err := postService.GetFeed(context.Background(), tt.userID, tt.limit, tt.offset)

			// Assert results
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
				assert.Len(t, posts, len(tt.expectedPosts))
				for i, expectedPost := range tt.expectedPosts {
					assert.Equal(t, expectedPost.Content, posts[i].Content)
				}
			}

			// Assert that all expectations were met
			mockPostRepo.AssertExpectations(t)
			mockUserRepo.AssertExpectations(t)
		})
	}
}
