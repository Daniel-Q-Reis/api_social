package integration_test

import (
	"context"
	"os"
	"testing"

	"social/api/internal/entity"
	repoInterface "social/api/internal/repo"
	"social/api/internal/repo/postgres"
	pg "social/api/pkg/postgres"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type PostgresInteractionTestSuite struct {
	suite.Suite
	db          *pg.Postgres
	userRepo    repoInterface.User
	postRepo    repoInterface.Post
	commentRepo repoInterface.Comment
	likeRepo    repoInterface.Like
	followRepo  repoInterface.Follow
	users       []*entity.User
	posts       []*entity.Post
}

// SetupSuite runs once before all tests in the suite
func (suite *PostgresInteraction_TestSuite) SetupSuite() {
	// Get database connection string from environment
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
	suite.commentRepo = postgres.NewCommentRepo(suite.db.Pool)
	suite.likeRepo = postgres.NewLikeRepo(suite.db.Pool)
	suite.followRepo = postgres.NewFollowRepo(suite.db.Pool)
}

// TearDownSuite runs once after all tests in the suite
func (suite *PostgresInteraction_TestSuite) TearDownSuite() {
	if suite.db != nil {
		suite.db.Close()
	}
}

// SetupTest runs before each test
func (suite *PostgresInteraction_TestSuite) SetupTest() {
	// Clean up database before each test
	queries := []string{
		"DELETE FROM likes",
		"DELETE FROM comments",
		"DELETE FROM posts",
		"DELETE FROM followers",
		"DELETE FROM users",
	}

	ctx := context.Background()
	for _, query := range queries {
		_, err := suite.db.Pool.Exec(ctx, query)
		if err != nil {
			suite.T().Fatalf("Failed to clean up database: %v", err)
		}
	}

	// Create test users
	suite.users = make([]*entity.User, 2)
	for i := 0; i < 2; i++ {
		suite.users[i] = &entity.User{
			Name:     "User " + string(rune('A'+i)),
			Username: "user" + string(rune('A'+i)),
			Email:    string(rune('A'+i)) + "@example.com",
			Password: "hashed_password",
		}
		err := suite.userRepo.Create(ctx, suite.users[i])
		if err != nil {
			suite.T().Fatalf("Failed to create user: %v", err)
		}
	}

	// Create a test post
	suite.posts = make([]*entity.Post, 1)
	suite.posts[0] = &entity.Post{
		AuthorID: suite.users[0].ID,
		Content:  "This is a test post",
	}
	err := suite.postRepo.Create(ctx, suite.posts[0])
	if err != nil {
		suite.T().Fatalf("Failed to create post: %v", err)
	}
}

// TestCommentRepo tests the Comment repository implementation
func (suite *PostgresInteraction_TestSuite) TestCommentRepo() {
	ctx := context.Background()

	// Test Create
	comment := &entity.Comment{
		PostID:   suite.posts[0].ID,
		AuthorID: suite.users[1].ID,
		Content:  "This is a test comment",
	}

	err := suite.commentRepo.Create(ctx, comment)
	assert.NoError(suite.T(), err)
	assert.NotEqual(suite.T(), uuid.Nil, comment.ID)

	// Test GetByPostID
	comments, err := suite.commentRepo.GetByPostID(ctx, suite.posts[0].ID, 10, 0)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), comments, 1)
	assert.Equal(suite.T(), comment.ID, comments[0].ID)
	assert.Equal(suite.T(), comment.PostID, comments[0].PostID)
	assert.Equal(suite.T(), comment.AuthorID, comments[0].AuthorID)
	assert.Equal(suite.T(), comment.Content, comments[0].Content)

	// Test Delete
	err = suite.commentRepo.Delete(ctx, comment.ID)
	assert.NoError(suite.T(), err)

	// Verify deletion
	comments, err = suite.commentRepo.GetByPostID(ctx, suite.posts[0].ID, 10, 0)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), comments, 0)
}

// TestLikeRepo tests the Like repository implementation
func (suite *PostgresInteraction_TestSuite) TestLikeRepo() {
	ctx := context.Background()

	// Test Create
	like := &entity.Like{
		UserID: suite.users[1].ID,
		PostID: suite.posts[0].ID,
	}

	err := suite.likeRepo.Create(ctx, like)
	assert.NoError(suite.T(), err)

	// Test Exists
	exists, err := suite.likeRepo.Exists(ctx, suite.users[1].ID, suite.posts[0].ID)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), exists)

	// Test Delete
	err = suite.likeRepo.Delete(ctx, suite.users[1].ID, suite.posts[0].ID)
	assert.NoError(suite.T(), err)

	// Verify deletion
	exists, err = suite.likeRepo.Exists(ctx, suite.users[1].ID, suite.posts[0].ID)
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), exists)
}

// TestFollowRepo tests the Follow repository implementation
func (suite *PostgresInteraction_TestSuite) TestFollowRepo() {
	ctx := context.Background()

	// Test Create
	follow := &entity.Follow{
		UserID:     suite.users[0].ID,
		FollowerID: suite.users[1].ID,
	}

	err := suite.followRepo.Create(ctx, follow)
	assert.NoError(suite.T(), err)

	// Test Exists
	exists, err := suite.followRepo.Exists(ctx, suite.users[0].ID, suite.users[1].ID)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), exists)

	// Test GetFollowers
	followers, err := suite.followRepo.GetFollowers(ctx, suite.users[0].ID, 10, 0)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), followers, 1)
	assert.Equal(suite.T(), suite.users[1].ID, followers[0].ID)

	// Test GetFollowing
	following, err := suite.followRepo.GetFollowing(ctx, suite.users[1].ID, 10, 0)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), following, 1)
	assert.Equal(suite.T(), suite.users[0].ID, following[0].ID)

	// Test Delete
	err = suite.followRepo.Delete(ctx, suite.users[0].ID, suite.users[1].ID)
	assert.NoError(suite.T(), err)

	// Verify deletion
	exists, err = suite.followRepo.Exists(ctx, suite.users[0].ID, suite.users[1].ID)
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), exists)
}

func TestPostgresInteractionSuite(t *testing.T) {
	// Skip integration tests if running in a CI environment without database
	if os.Getenv("SKIP_INTEGRATION_TESTS") == "true" {
		t.Skip("Skipping integration tests")
	}

	suite.Run(t, new(PostgresInteraction_TestSuite))
}
