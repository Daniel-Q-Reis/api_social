package postgres

import (
	"context"
	"fmt"

	"social/api/internal/entity"
	"social/api/internal/repo"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type LikeRepo struct {
	db *pgxpool.Pool
}

func NewLikeRepo(db *pgxpool.Pool) repo.Like {
	return &LikeRepo{db: db}
}

func (r *LikeRepo) Create(ctx context.Context, like *entity.Like) error {
	query := `INSERT INTO likes (user_id, post_id) VALUES ($1, $2) ON CONFLICT DO NOTHING RETURNING created_at`
	err := r.db.QueryRow(ctx, query, like.UserID, like.PostID).Scan(&like.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create like: %w", err)
	}

	return nil
}

func (r *LikeRepo) Delete(ctx context.Context, userID, postID uuid.UUID) error {
	query := `DELETE FROM likes WHERE user_id = $1 AND post_id = $2`
	result, err := r.db.Exec(ctx, query, userID, postID)

	if err != nil {
		return fmt.Errorf("failed to delete like: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("like not found")
	}

	return nil
}

func (r *LikeRepo) Exists(ctx context.Context, userID, postID uuid.UUID) (bool, error) {
	var (
		exists bool
		query  = `SELECT EXISTS(SELECT 1 FROM likes WHERE user_id = $1 AND post_id = $2)`
	)

	err := r.db.QueryRow(ctx, query, userID, postID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check if like exists: %w", err)
	}

	return exists, nil
}
