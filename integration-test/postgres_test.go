package integration_test

import (
	"context"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"social/api/internal/entity"
	repoInterface "social/api/internal/repo"
	"social/api/internal/repo/postgres"
	pg "social/api/pkg/postgres"
)

type Postgres_TestSuite struct {
	suite.Suite
	db       *pg.Postgres
	userRepo repoInterface.User
	postRepo repoInterface.Post
}

// SetupSuite runs once before all tests in the suite
func (suite *Postgres_TestSuite) SetupSuite() {
	// Get database connection string from environment
	// In integration tests, this would typically be set to point to the test database
	dbURL := os.Getenv("PG_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:password@localhost:5432/social_db?sslmode=disable"
	}

	// Connect to the database
	var err error
	suite.db, err = pg.New(dbURL)
	if err != nil {
		suite.T().Fatalf("Failed to connect to database: %v", err)
	}

	// Initialize repositories
	suite.userRepo = postgres.NewUserRepo(suite.db.Pool)
	suite.postRepo = postgres.NewPostRepo(suite.db.Pool)
}

// TearDownSuite runs once after all tests in the suite
func (suite *Postgres_TestSuite) TearDownSuite() {
	if suite.db != nil {
		suite.db.Close()
	}
}

// SetupTest runs before each test
func (suite *Postgres_TestSuite) SetupTest() {
	// Clean up database before each test
	queries := []string{
		"DELETE FROM likes",
		"DELETE FROM comments",
		"DELETE FROM posts",
		"DELETE FROM followers",
		"DELETE FROM users",
	}
	
	for _, query := range queries {
		_, err := suite.db.Pool.Exec(context.Background(), query)
		if err != nil {
			suite.T().Fatalf("Failed to clean up database: %v", err)
		}
	}
}

// TestUserRepo tests the User repository implementation
func (suite *Postgres_TestSuite) TestUserRepo() {
	// Test Create
	user := &entity.User{
		Name:     "John Doe",
		Username: "johndoe",
		Email:    "john@example.com",
		Password: "hashed_password",
		Bio:      stringPtr("This is my bio"),
		ImageURL: stringPtr("https://example.com/image.jpg"),
	}
	
	err := suite.userRepo.Create(context.Background(), user)
	assert.NoError(suite.T(), err)
	assert.NotEqual(suite.T(), uuid.Nil, user.ID)
	
	// Test GetByID
	retrievedUser, err := suite.userRepo.GetByID(context.Background(), user.ID)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), user.ID, retrievedUser.ID)
	assert.Equal(suite.T(), user.Name, retrievedUser.Name)
	assert.Equal(suite.T(), user.Username, retrievedUser.Username)
	assert.Equal(suite.T(), user.Email, retrievedUser.Email)
	assert.Equal(suite.T(), user.Password, retrievedUser.Password)
	assert.Equal(suite.T(), *user.Bio, *retrievedUser.Bio)
	assert.Equal(suite.T(), *user.ImageURL, *retrievedUser.ImageURL)
	
	// Test GetByEmail
	retrievedUser, err = suite.userRepo.GetByEmail(context.Background(), user.Email)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), user.Email, retrievedUser.Email)
	
	// Test GetByUsername
	retrievedUser, err = suite.userRepo.GetByUsername(context.Background(), user.Username)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), user.Username, retrievedUser.Username)
	
	// Test Update
	newName := "Jane Doe"
	newBio := "Updated bio"
	user.Name = newName
	user.Bio = &newBio
	
	err = suite.userRepo.Update(context.Background(), user)
	assert.NoError(suite.T(), err)
	
	// Verify update
	updatedUser, err := suite.userRepo.GetByID(context.Background(), user.ID)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), newName, updatedUser.Name)
	assert.Equal(suite.T(), newBio, *updatedUser.Bio)
	
	// Test Search
	users, err := suite.userRepo.Search(context.Background(), "joh", 10, 0)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), users, 1)
	assert.Equal(suite.T(), user.Username, users[0].Username)
	
	// Test duplicate email
	duplicateUser := &entity.User{
		Name:     "Jane Smith",
		Username: "janesmith",
		Email:    "john@example.com", // Same email
		Password: "hashed_password",
	}
	
	err = suite.userRepo.Create(context.Background(), duplicateUser)
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "user already exists")
	
	// Test duplicate username
	duplicateUser2 := &entity.User{
		Name:     "Jane Smith",
		Username: "johndoe", // Same username
		Email:    "jane@example.com",
		Password: "hashed_password",
	}
	
	err = suite.userRepo.Create(context.Background(), duplicateUser2)
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "user already exists")
}

// TestPostRepo tests the Post repository implementation
func (suite *Postgres_TestSuite) TestPostRepo() {
	// First create a user for the post
	user := &entity.User{
		Name:     "John Doe",
		Username: "johndoe",
		Email:    "john@example.com",
		Password: "hashed_password",
	}
	
	err := suite.userRepo.Create(context.Background(), user)
	assert.NoError(suite.T(), err)
	
	// Test Create
	post := &entity.Post{
		AuthorID: user.ID,
		Content:  "This is a test post",
		ImageURL: stringPtr("https://example.com/post-image.jpg"),
	}
	
	err = suite.postRepo.Create(context.Background(), post)
	assert.NoError(suite.T(), err)
	assert.NotEqual(suite.T(), uuid.Nil, post.ID)
	
	// Test GetByID
	retrievedPost, err := suite.postRepo.GetByID(context.Background(), post.ID)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), post.ID, retrievedPost.ID)
	assert.Equal(suite.T(), post.AuthorID, retrievedPost.AuthorID)
	assert.Equal(suite.T(), post.Content, retrievedPost.Content)
	assert.Equal(suite.T(), *post.ImageURL, *retrievedPost.ImageURL)
	
	// Test GetByAuthorID
	posts, err := suite.postRepo.GetByAuthorID(context.Background(), user.ID, 10, 0)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), posts, 1)
	assert.Equal(suite.T(), post.ID, posts[0].ID)
	
	// Test Update
	newContent := "Updated post content"
	post.Content = newContent
	
	err = suite.postRepo.Update(context.Background(), post)
	assert.NoError(suite.T(), err)
	
	// Verify update
	updatedPost, err := suite.postRepo.GetByID(context.Background(), post.ID)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), newContent, updatedPost.Content)
	
	// Test Delete
	err = suite.postRepo.Delete(context.Background(), post.ID)
	assert.NoError(suite.T(), err)
	
	// Verify deletion
	_, err = suite.postRepo.GetByID(context.Background(), post.ID)
	assert.Error(suite.T(), err)
}

func TestPostgresSuite(t *testing.T) {
	// Skip integration tests if running in a CI environment without database
	if os.Getenv("SKIP_INTEGRATION_TESTS") == "true" {
		t.Skip("Skipping integration tests")
	}
	
	suite.Run(t, new(Postgres_TestSuite))
}

// Helper function to create string pointer
func stringPtr(s string) *string {
	return &s
}