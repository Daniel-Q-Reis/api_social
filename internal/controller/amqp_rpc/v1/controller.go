package v1

import (
	"github.com/Daniel-Q-Reis/ia_social_media/internal/usecase"
	"github.com/Daniel-Q-Reis/ia_social_media/pkg/logger"
	"github.com/go-playground/validator/v10"
)

// V1 -.
type V1 struct {
	t usecase.Translation
	l logger.Interface
	v *validator.Validate
}
