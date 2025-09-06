package entity

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type Comment struct {
	ID        uuid.UUID `json:"id" db:"id"`
	PostID    uuid.UUID `json:"post_id" db:"post_id"`
	AuthorID  uuid.UUID `json:"author_id" db:"author_id"`
	Content   string    `json:"content" db:"content" validate:"required,min=1,max=500"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

func (c *Comment) Validate() error {
	validate := validator.New()
	return validate.Struct(c)
}