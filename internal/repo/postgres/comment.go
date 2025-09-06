package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"social/api/internal/entity"
	"social/api/internal/repo"
)

type CommentRepo struct {
	db *pgxpool.Pool
}

func NewCommentRepo(db *pgxpool.Pool) repo.Comment {
	return &CommentRepo{db: db}
}

func (r *CommentRepo) Create(ctx context.Context, comment *entity.Comment) error {
	query := `INSERT INTO comments (post_id, author_id, content) 
	          VALUES ($1, $2, $3) RETURNING id, created_at`
	err := r.db.QueryRow(ctx, query, comment.PostID, comment.AuthorID, comment.Content).Scan(&comment.ID, &comment.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create comment: %w", err)
	}
	return nil
}

func (r *CommentRepo) GetByPostID(ctx context.Context, postID uuid.UUID) ([]entity.Comment, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, post_id, author_id, content, created_at 
		FROM comments 
		WHERE post_id = $1 
		ORDER BY created_at ASC`, postID)
	if err != nil {
		return nil, fmt.Errorf("failed to get comments by post ID: %w", err)
	}
	defer rows.Close()

	var comments []entity.Comment
	for rows.Next() {
		var comment entity.Comment
		err := rows.Scan(&comment.ID, &comment.PostID, &comment.AuthorID, &comment.Content, &comment.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan comment: %w", err)
		}
		comments = append(comments, comment)
	}

	return comments, nil
}

func (r *CommentRepo) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM comments WHERE id = $1`
	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete comment: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("comment not found")
	}
	return nil
}