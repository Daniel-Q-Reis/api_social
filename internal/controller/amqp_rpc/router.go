package v1

import (
	v1 "github.com/Daniel-Q-Reis/ia_social_media/internal/controller/amqp_rpc/v1"
	"github.com/Daniel-Q-Reis/ia_social_media/internal/usecase"
	"github.com/Daniel-Q-Reis/ia_social_media/pkg/logger"
	"github.com/Daniel-Q-Reis/ia_social_media/pkg/rabbitmq/rmq_rpc/server"
)

// NewRouter -.
func NewRouter(t usecase.Translation, l logger.Interface) map[string]server.CallHandler {
	routes := make(map[string]server.CallHandler)

	{
		v1.NewTranslationRoutes(routes, t, l)
	}

	return routes
}
