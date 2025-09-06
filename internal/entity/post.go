package entity

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type Post struct {
	ID        uuid.UUID `json:"id" db:"id"`
	AuthorID  uuid.UUID `json:"author_id" db:"author_id"`
	Content   string    `json:"content" db:"content" validate:"required,min=1,max=1000"`
	ImageURL  *string   `json:"image_url,omitempty" db:"image_url" validate:"omitempty,url"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

func (p *Post) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}