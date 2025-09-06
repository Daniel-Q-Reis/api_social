package entity

import (
	"time"

	"github.com/google/uuid"
)

type Like struct {
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	PostID    uuid.UUID `json:"post_id" db:"post_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}