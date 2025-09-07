package entity

import (
	"time"

	"github.com/google/uuid"
)

type Follow struct {
	UserID     uuid.UUID `json:"user_id" db:"user_id"`
	FollowerID uuid.UUID `json:"follower_id" db:"follower_id"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}
