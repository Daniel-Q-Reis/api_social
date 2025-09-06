package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"social/api/internal/entity"
	"social/api/internal/repo"
)

type PostRepo struct {
	db *pgxpool.Pool
}

func NewPostRepo(db *pgxpool.Pool) repo.Post {
	return &PostRepo{db: db}
}

func (r *PostRepo) Create(ctx context.Context, post *entity.Post) error {
	query := `INSERT INTO posts (author_id, content, image_url) 
	          VALUES ($1, $2, $3) RETURNING id, created_at, updated_at`
	err := r.db.QueryRow(ctx, query, post.AuthorID, post.Content, post.ImageURL).Scan(&post.ID, &post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create post: %w", err)
	}
	return nil
}

func (r *PostRepo) GetByID(ctx context.Context, id uuid.UUID) (*entity.Post, error) {
	var post entity.Post
	query := `SELECT id, author_id, content, image_url, created_at, updated_at 
	          FROM posts WHERE id = $1`
	err := r.db.QueryRow(ctx, query, id).Scan(
		&post.ID, &post.AuthorID, &post.Content, &post.ImageURL, &post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repo.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get post by ID: %w", err)
	}
	return &post, nil
}

func (r *PostRepo) GetByAuthorID(ctx context.Context, authorID uuid.UUID, limit, offset int) ([]entity.Post, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, author_id, content, image_url, created_at, updated_at 
		FROM posts 
		WHERE author_id = $1 
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`, authorID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get posts by author ID: %w", err)
	}
	defer rows.Close()

	var posts []entity.Post
	for rows.Next() {
		var post entity.Post
		err := rows.Scan(&post.ID, &post.AuthorID, &post.Content, &post.ImageURL, &post.CreatedAt, &post.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan post: %w", err)
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func (r *PostRepo) GetFeed(ctx context.Context, userID uuid.UUID, limit, offset int) ([]entity.Post, error) {
	rows, err := r.db.Query(ctx, `
		SELECT p.id, p.author_id, p.content, p.image_url, p.created_at, p.updated_at
		FROM posts p
		JOIN followers f ON p.author_id = f.user_id
		WHERE f.follower_id = $1
		ORDER BY p.created_at DESC
		LIMIT $2 OFFSET $3`, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get feed: %w", err)
	}
	defer rows.Close()

	var posts []entity.Post
	for rows.Next() {
		var post entity.Post
		err := rows.Scan(&post.ID, &post.AuthorID, &post.Content, &post.ImageURL, &post.CreatedAt, &post.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan post: %w", err)
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func (r *PostRepo) Update(ctx context.Context, post *entity.Post) error {
	query := `UPDATE posts SET content = $1, image_url = $2, updated_at = NOW() 
	          WHERE id = $3 RETURNING updated_at`
	err := r.db.QueryRow(ctx, query, post.Content, post.ImageURL, post.ID).Scan(&post.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return repo.ErrNotFound
		}
		return fmt.Errorf("failed to update post: %w", err)
	}
	return nil
}

func (r *PostRepo) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM posts WHERE id = $1`
	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete post: %w", err)
	}
	if result.RowsAffected() == 0 {
		return repo.ErrNotFound
	}
	return nil
}