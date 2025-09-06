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

func TestCommentService_AddComment(t *testing.T) {
	postID := uuid.New()
	userID := uuid.New()
	content := "This is a test comment"

	// Test cases
	tests := []struct {
		name          string
		postID        uuid.UUID
		userID        uuid.UUID
		content       string
		setupMock     func(*CommentRepoMock, *UserRepoMock, *PostRepoMock)
		expectedComment *entity.Comment
		expectedError   error
	}{
		{
			name:    "Success",
			postID:  postID,
			userID:  userID,
			content: content,
			setupMock: func(mockCommentRepo *CommentRepoMock, mockUserRepo *UserRepoMock, mockPostRepo *PostRepoMock) {
				// Mock GetByID to return a post (verify post exists)
				post := &entity.Post{
					ID:       postID,
					AuthorID: uuid.New(),
					Content:  "Test post",
				}
				mockPostRepo.On("GetByID", mock.Anything, postID).Return(post, nil)

				// Mock Create to succeed
				mockCommentRepo.On("Create", mock.Anything, mock.MatchedBy(func(comment *entity.Comment) bool {
					return comment.PostID == postID && comment.AuthorID == userID && comment.Content == content
				})).Return(nil)
			},
			expectedComment: &entity.Comment{
				PostID:   postID,
				AuthorID: userID,
				Content:  content,
			},
			expectedError: nil,
		},
		{
			name:    "PostNotFound",
			postID:  postID,
			userID:  userID,
			content: content,
			setupMock: func(mockCommentRepo *CommentRepoMock, mockUserRepo *UserRepoMock, mockPostRepo *PostRepoMock) {
				// Mock GetByID to return ErrNotFound
				mockPostRepo.On("GetByID", mock.Anything, postID).Return(nil, repo.ErrNotFound)
			},
			expectedComment: nil,
			expectedError:   repo.ErrNotFound,
		},
		{
			name:    "ValidationFailed_EmptyContent",
			postID:  postID,
			userID:  userID,
			content: "", // Empty content
			setupMock: func(mockCommentRepo *CommentRepoMock, mockUserRepo *UserRepoMock, mockPostRepo *PostRepoMock) {
				// Mock GetByID to return a post (verify post exists)
				post := &entity.Post{
					ID:       postID,
					AuthorID: uuid.New(),
					Content:  "Test post",
				}
				mockPostRepo.On("GetByID", mock.Anything, postID).Return(post, nil)
			},
			expectedComment: nil,
			expectedError:   errors.New("validation failed"),
		},
		{
			name:    "DatabaseErrorOnCreate",
			postID:  postID,
			userID:  userID,
			content: content,
			setupMock: func(mockCommentRepo *CommentRepoMock, mockUserRepo *UserRepoMock, mockPostRepo *PostRepoMock) {
				// Mock GetByID to return a post (verify post exists)
				post := &entity.Post{
					ID:       postID,
					AuthorID: uuid.New(),
					Content:  "Test post",
				}
				mockPostRepo.On("GetByID", mock.Anything, postID).Return(post, nil)

				// Mock Create to return an error
				mockCommentRepo.On("Create", mock.Anything, mock.Anything).Return(errors.New("database error"))
			},
			expectedComment: nil,
			expectedError:   errors.New("failed to add comment"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock repositories
			mockCommentRepo := &CommentRepoMock{}
			mockUserRepo := &UserRepoMock{}
			mockPostRepo := &PostRepoMock{}

			// Setup mock expectations
			tt.setupMock(mockCommentRepo, mockUserRepo, mockPostRepo)

			// Create comment service with mock repositories
			commentService := usecase.NewCommentUseCase(mockCommentRepo, mockUserRepo, mockPostRepo)

			// Execute the method under test
			comment, err := commentService.AddComment(context.Background(), tt.postID, tt.userID, tt.content)

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
				assert.NotNil(t, comment)
				assert.Equal(t, tt.expectedComment.PostID, comment.PostID)
				assert.Equal(t, tt.expectedComment.AuthorID, comment.AuthorID)
				assert.Equal(t, tt.expectedComment.Content, comment.Content)
				// Verify that ID was set (not zero value)
				assert.NotEqual(t, uuid.Nil, comment.ID)
			}

			// Assert that all expectations were met
			mockCommentRepo.AssertExpectations(t)
			mockUserRepo.AssertExpectations(t)
			mockPostRepo.AssertExpectations(t)
		})
	}
}

func TestCommentService_GetComments(t *testing.T) {
	postID := uuid.New()
	commentID1 := uuid.New()
	commentID2 := uuid.New()
	userID := uuid.New()

	// Test cases
	tests := []struct {
		name            string
		postID          uuid.UUID
		limit           int
		offset          int
		setupMock       func(*CommentRepoMock, *UserRepoMock, *PostRepoMock)
		expectedComments []entity.Comment
		expectedError    error
	}{
		{
			name:   "Success",
			postID: postID,
			limit:  10,
			offset: 0,
			setupMock: func(mockCommentRepo *CommentRepoMock, mockUserRepo *UserRepoMock, mockPostRepo *PostRepoMock) {
				// Mock GetByID to return a post (verify post exists)
				post := &entity.Post{
					ID:       postID,
					AuthorID: uuid.New(),
					Content:  "Test post",
				}
				mockPostRepo.On("GetByID", mock.Anything, postID).Return(post, nil)

				// Mock GetByPostID to return comments
				comments := []entity.Comment{
					{
						ID:       commentID1,
						PostID:   postID,
						AuthorID: userID,
						Content:  "First comment",
					},
					{
						ID:       commentID2,
						PostID:   postID,
						AuthorID: userID,
						Content:  "Second comment",
					},
				}
				mockCommentRepo.On("GetByPostID", mock.Anything, postID, 10, 0).Return(comments, nil)
			},
			expectedComments: []entity.Comment{
				{
					ID:       commentID1,
					PostID:   postID,
					AuthorID: userID,
					Content:  "First comment",
				},
				{
					ID:       commentID2,
					PostID:   postID,
					AuthorID: userID,
					Content:  "Second comment",
				},
			},
			expectedError: nil,
		},
		{
			name:   "PostNotFound",
			postID: postID,
			limit:  10,
			offset: 0,
			setupMock: func(mockCommentRepo *CommentRepoMock, mockUserRepo *UserRepoMock, mockPostRepo *PostRepoMock) {
				// Mock GetByID to return ErrNotFound
				mockPostRepo.On("GetByID", mock.Anything, postID).Return(nil, repo.ErrNotFound)
			},
			expectedComments: nil,
			expectedError:    repo.ErrNotFound,
		},
		{
			name:   "NoComments",
			postID: postID,
			limit:  10,
			offset: 0,
			setupMock: func(mockCommentRepo *CommentRepoMock, mockUserRepo *UserRepoMock, mockPostRepo *PostRepoMock) {
				// Mock GetByID to return a post (verify post exists)
				post := &entity.Post{
					ID:       postID,
					AuthorID: uuid.New(),
					Content:  "Test post",
				}
				mockPostRepo.On("GetByID", mock.Anything, postID).Return(post, nil)

				// Mock GetByPostID to return empty slice
				mockCommentRepo.On("GetByPostID", mock.Anything, postID, 10, 0).Return([]entity.Comment{}, nil)
			},
			expectedComments: []entity.Comment{},
			expectedError:    nil,
		},
		{
			name:   "DatabaseErrorOnGetComments",
			postID: postID,
			limit:  10,
			offset: 0,
			setupMock: func(mockCommentRepo *CommentRepoMock, mockUserRepo *UserRepoMock, mockPostRepo *PostRepoMock) {
				// Mock GetByID to return a post (verify post exists)
				post := &entity.Post{
					ID:       postID,
					AuthorID: uuid.New(),
					Content:  "Test post",
				}
				mockPostRepo.On("GetByID", mock.Anything, postID).Return(post, nil)

				// Mock GetByPostID to return an error
				mockCommentRepo.On("GetByPostID", mock.Anything, postID, 10, 0).Return(nil, errors.New("database error"))
			},
			expectedComments: nil,
			expectedError:    errors.New("failed to get comments"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock repositories
			mockCommentRepo := &CommentRepoMock{}
			mockUserRepo := &UserRepoMock{}
			mockPostRepo := &PostRepoMock{}

			// Setup mock expectations
			tt.setupMock(mockCommentRepo, mockUserRepo, mockPostRepo)

			// Create comment service with mock repositories
			commentService := usecase.NewCommentUseCase(mockCommentRepo, mockUserRepo, mockPostRepo)

			// Execute the method under test
			comments, err := commentService.GetComments(context.Background(), tt.postID, tt.limit, tt.offset)

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
				assert.Len(t, comments, len(tt.expectedComments))
				for i, expectedComment := range tt.expectedComments {
					assert.Equal(t, expectedComment.ID, comments[i].ID)
					assert.Equal(t, expectedComment.PostID, comments[i].PostID)
					assert.Equal(t, expectedComment.AuthorID, comments[i].AuthorID)
					assert.Equal(t, expectedComment.Content, comments[i].Content)
				}
			}

			// Assert that all expectations were met
			mockCommentRepo.AssertExpectations(t)
			mockUserRepo.AssertExpectations(t)
			mockPostRepo.AssertExpectations(t)
		})
	}
}

func TestCommentService_DeleteComment(t *testing.T) {
	commentID := uuid.New()

	// Test cases
	tests := []struct {
		name          string
		commentID     uuid.UUID
		userID        uuid.UUID
		setupMock     func(*CommentRepoMock, *UserRepoMock, *PostRepoMock)
		expectedError error
	}{
		{
			name:      "Success",
			commentID: commentID,
			userID:    uuid.New(),
			setupMock: func(mockCommentRepo *CommentRepoMock, mockUserRepo *UserRepoMock, mockPostRepo *PostRepoMock) {
				// Mock Delete to succeed
				mockCommentRepo.On("Delete", mock.Anything, commentID).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:      "CommentNotFound",
			commentID: commentID,
			userID:    uuid.New(),
			setupMock: func(mockCommentRepo *CommentRepoMock, mockUserRepo *UserRepoMock, mockPostRepo *PostRepoMock) {
				// Mock Delete to return ErrNotFound
				mockCommentRepo.On("Delete", mock.Anything, commentID).Return(repo.ErrNotFound)
			},
			expectedError: repo.ErrNotFound,
		},
		{
			name:      "DatabaseError",
			commentID: commentID,
			userID:    uuid.New(),
			setupMock: func(mockCommentRepo *CommentRepoMock, mockUserRepo *UserRepoMock, mockPostRepo *PostRepoMock) {
				// Mock Delete to return an error
				mockCommentRepo.On("Delete", mock.Anything, commentID).Return(errors.New("database error"))
			},
			expectedError: errors.New("failed to delete comment"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock repositories
			mockCommentRepo := &CommentRepoMock{}
			mockUserRepo := &UserRepoMock{}
			mockPostRepo := &PostRepoMock{}

			// Setup mock expectations
			tt.setupMock(mockCommentRepo, mockUserRepo, mockPostRepo)

			// Create comment service with mock repositories
			commentService := usecase.NewCommentUseCase(mockCommentRepo, mockUserRepo, mockPostRepo)

			// Execute the method under test
			err := commentService.DeleteComment(context.Background(), tt.commentID, tt.userID)

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
			}

			// Assert that all expectations were met
			mockCommentRepo.AssertExpectations(t)
			mockUserRepo.AssertExpectations(t)
			mockPostRepo.AssertExpectations(t)
		})
	}
}