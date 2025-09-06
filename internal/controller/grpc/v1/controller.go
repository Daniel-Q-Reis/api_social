package v1

import (
	v1 "github.com/Daniel-Q-Reis/ia_social_media/docs/proto/v1"
	"github.com/Daniel-Q-Reis/ia_social_media/internal/usecase"
	"github.com/Daniel-Q-Reis/ia_social_media/pkg/logger"
	"github.com/go-playground/validator/v10"
)

// V1 -.
type V1 struct {
	v1.TranslationServer

	t usecase.Translation
	l logger.Interface
	v *validator.Validate
}
