package entity

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Name      string    `json:"name" db:"name" validate:"required,min=1,max=100"`
	Username  string    `json:"username" db:"username" validate:"required,min=3,max=50,alphanum"`
	Email     string    `json:"email" db:"email" validate:"required,email"`
	Password  string    `json:"-" db:"password_hash" validate:"required,min=6"`
	Bio       *string   `json:"bio,omitempty" db:"bio" validate:"omitempty,max=500"`
	ImageURL  *string   `json:"image_url,omitempty" db:"profile_picture_url" validate:"omitempty,url"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

func (u *User) Validate() error {
	validate := validator.New()
	return validate.Struct(u)
}