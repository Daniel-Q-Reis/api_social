package postgres

import (
	"context"
	"errors"
	"fmt"

	"social/api/internal/entity"
	"social/api/internal/repo"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepo struct {
	db *pgxpool.Pool
}

func NewUserRepo(db *pgxpool.Pool) repo.User {
	return &UserRepo{db: db}
}

func (r *UserRepo) Create(ctx context.Context, user *entity.User) error {
	query := `INSERT INTO users (name, username, email, password_hash, bio, profile_picture_url) 
	          VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created_at, updated_at`
	err := r.db.QueryRow(ctx, query, user.Name, user.Username, user.Email, user.Password, user.Bio, user.ImageURL).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		// Check for duplicate username or email
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			// Check for unique constraint violation
			if pgErr.Code == "23505" { // unique_violation
				// We need to check which constraint was violated
				// This is a simplified approach - in production, you'd check the specific constraint name
				return fmt.Errorf("user already exists: %w", err)
			}
		}

		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (r *UserRepo) GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	return r.getUserByField(ctx, "id", id)
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	return r.getUserByField(ctx, "email", email)
}

func (r *UserRepo) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	return r.getUserByField(ctx, "username", username)
}

func (r *UserRepo) getUserByField(ctx context.Context, field string, value interface{}) (*entity.User, error) {
	var (
		user  entity.User
		query = fmt.Sprintf(`SELECT id, name, username, email, password_hash, bio, profile_picture_url, created_at, updated_at 
	          FROM users WHERE %s = $1`, field)
	)

	err := r.db.QueryRow(ctx, query, value).Scan(
		&user.ID, &user.Name, &user.Username, &user.Email, &user.Password,
		&user.Bio, &user.ImageURL, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repo.ErrNotFound
		}

		return nil, fmt.Errorf("failed to get user by %s: %w", field, err)
	}

	return &user, nil
}

func (r *UserRepo) Update(ctx context.Context, user *entity.User) error {
	query := `UPDATE users SET name = $1, bio = $2, profile_picture_url = $3, updated_at = NOW() 
	          WHERE id = $4 RETURNING updated_at`
	err := r.db.QueryRow(ctx, query, user.Name, user.Bio, user.ImageURL, user.ID).Scan(&user.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return repo.ErrNotFound
		}

		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

func (r *UserRepo) Search(ctx context.Context, query string, limit, offset int) ([]entity.User, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, name, username, email, password_hash, bio, profile_picture_url, created_at, updated_at 
		FROM users 
		WHERE name ILIKE $1 OR username ILIKE $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`, "%"+query+"%", limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to search users: %w", err)
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
