package v1

import (
	"github.com/Daniel-Q-Reis/ia_social_media/internal/usecase"
	"github.com/Daniel-Q-Reis/ia_social_media/pkg/logger"
	"github.com/Daniel-Q-Reis/ia_social_media/pkg/rabbitmq/rmq_rpc/server"
	"github.com/go-playground/validator/v10"
)

// NewTranslationRoutes -.
func NewTranslationRoutes(routes map[string]server.CallHandler, t usecase.Translation, l logger.Interface) {
	r := &V1{t: t, l: l, v: validator.New(validator.WithRequiredStructEnabled())}

	{
		routes["v1.getHistory"] = r.getHistory()
	}
}
