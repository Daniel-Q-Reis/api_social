package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"social/api/internal/entity"
	"social/api/internal/repo"
)

type FollowRepo struct {
	db *pgxpool.Pool
}

func NewFollowRepo(db *pgxpool.Pool) repo.Follow {
	return &FollowRepo{db: db}
}

func (r *FollowRepo) Create(ctx context.Context, follow *entity.Follow) error {
	query := `INSERT INTO followers (user_id, follower_id) VALUES ($1, $2) ON CONFLICT DO NOTHING RETURNING created_at`
	err := r.db.QueryRow(ctx, query, follow.UserID, follow.FollowerID).Scan(&follow.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create follow: %w", err)
	}
	return nil
}

func (r *FollowRepo) Delete(ctx context.Context, userID, followerID uuid.UUID) error {
	query := `DELETE FROM followers WHERE user_id = $1 AND follower_id = $2`
	result, err := r.db.Exec(ctx, query, userID, followerID)
	if err != nil {
		return fmt.Errorf("failed to delete follow: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("follow relationship not found")
	}
	return nil
}

func (r *FollowRepo) Exists(ctx context.Context, userID, followerID uuid.UUID) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM followers WHERE user_id = $1 AND follower_id = $2)`
	err := r.db.QueryRow(ctx, query, userID, followerID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check if follow exists: %w", err)
	}
	return exists, nil
}

func (r *FollowRepo) GetFollowers(ctx context.Context, userID uuid.UUID) ([]entity.User, error) {
	rows, err := r.db.Query(ctx, `
		SELECT u.id, u.name, u.username, u.email, u.password_hash, u.bio, u.profile_picture_url, u.created_at, u.updated_at
		FROM users u
		JOIN followers f ON u.id = f.follower_id
		WHERE f.user_id = $1
		ORDER BY f.created_at DESC`, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get followers: %w", err)
	}
	defer rows.Close()

	var users []entity.User
	for rows.Next() {
		var user entity.User
		err := rows.Scan(
			&user.ID, &user.Name, &user.Username, &user.Email, &user.Password,
			&user.Bio, &user.ImageURL, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *FollowRepo) GetFollowing(ctx context.Context, followerID uuid.UUID) ([]entity.User, error) {
	rows, err := r.db.Query(ctx, `
		SELECT u.id, u.name, u.username, u.email, u.password_hash, u.bio, u.profile_picture_url, u.created_at, u.updated_at
		FROM users u
		JOIN followers f ON u.id = f.user_id
		WHERE f.follower_id = $1
		ORDER BY f.created_at DESC`, followerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get following: %w", err)
	}
	defer rows.Close()

	var users []entity.User
	for rows.Next() {
		var user entity.User
		err := rows.Scan(
			&user.ID, &user.Name, &user.Username, &user.Email, &user.Password,
			&user.Bio, &user.ImageURL, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	return users, nil
}