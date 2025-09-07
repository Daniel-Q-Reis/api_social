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

func TestCommentService_AddComment(t *testing.T) {
	tests := []struct {
		name            string
		mockUserRepo    *UserRepoMock
		mockCommentRepo *CommentRepoMock
		postID          uuid.UUID
		userID          uuid.UUID
		content         string
		wantErr         bool
	}{
		{
			name:            "ValidComment",
			mockUserRepo:    &UserRepoMock{},
			mockCommentRepo: &CommentRepoMock{},
			postID:          uuid.New(),
			userID:          uuid.New(),
			content:         "This is a test comment",
			wantErr:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock expectations
			if !tt.wantErr {
				// For AddComment, we only need to mock the post existence check
				tt.mockCommentRepo.On("Create", mock.Anything, mock.MatchedBy(func(c *entity.Comment) bool {
					return c.PostID == tt.postID && c.AuthorID == tt.userID && c.Content == tt.content
				})).Return(nil)
			}

			// We need to create a mock PostRepo for the AddComment method
			mockPostRepo := &PostRepoMock{}
			// Mock the post existence check
			mockPostRepo.On("GetByID", mock.Anything, tt.postID).Return(&entity.Post{ID: tt.postID}, nil)

			s := usecase.NewCommentUseCase(tt.mockCommentRepo, tt.mockUserRepo, mockPostRepo)
			got, err := s.AddComment(context.Background(), tt.postID, tt.userID, tt.content)

			if (err != nil) != tt.wantErr {
				t.Errorf("CommentService.AddComment() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				// For a successful test, we just check that we got a comment back with the right properties
				if got == nil {
					t.Errorf("CommentService.AddComment() = nil, want comment")
					return
				}
				if got.PostID != tt.postID {
					t.Errorf("CommentService.AddComment() PostID = %v, want %v", got.PostID, tt.postID)
				}

				if got.AuthorID != tt.userID {
					t.Errorf("CommentService.AddComment() AuthorID = %v, want %v", got.AuthorID, tt.userID)
				}
				if got.Content != tt.content {
					t.Errorf("CommentService.AddComment() Content = %v, want %v", got.Content, tt.content)
				}

				if got.ID == uuid.Nil {
					t.Errorf("CommentService.AddComment() ID should not be nil")
				}
			}

			// Assert that all expectations were met
			tt.mockUserRepo.AssertExpectations(t)
			tt.mockCommentRepo.AssertExpectations(t)
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
		name             string
		postID           uuid.UUID
		limit            int
		offset           int
		setupMock        func(*CommentRepoMock, *UserRepoMock, *PostRepoMock)
		expectedComments []entity.Comment
		expectedError    error
	}{
		{
			name:   "Success",
			postID: postID,
			limit:  10,
			offset: 0,
			setupMock: func(mockCommentRepo *CommentRepoMock, _ *UserRepoMock, mockPostRepo *PostRepoMock) {
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
			setupMock: func(_ *CommentRepoMock, _ *UserRepoMock, mockPostRepo *PostRepoMock) {
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
			setupMock: func(mockCommentRepo *CommentRepoMock, _ *UserRepoMock, mockPostRepo *PostRepoMock) {
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
			setupMock: func(mockCommentRepo *CommentRepoMock, _ *UserRepoMock, mockPostRepo *PostRepoMock) {
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
			setupMock: func(mockCommentRepo *CommentRepoMock, _ *UserRepoMock, _ *PostRepoMock) {
				// Mock Delete to succeed
				mockCommentRepo.On("Delete", mock.Anything, commentID).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:      "CommentNotFound",
			commentID: commentID,
			userID:    uuid.New(),
			setupMock: func(mockCommentRepo *CommentRepoMock, _ *UserRepoMock, _ *PostRepoMock) {
				// Mock Delete to return ErrNotFound
				mockCommentRepo.On("Delete", mock.Anything, commentID).Return(repo.ErrNotFound)
			},
			expectedError: repo.ErrNotFound,
		},
		{
			name:      "DatabaseError",
			commentID: commentID,
			userID:    uuid.New(),
			setupMock: func(mockCommentRepo *CommentRepoMock, _ *UserRepoMock, _ *PostRepoMock) {
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
