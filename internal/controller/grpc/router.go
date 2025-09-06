package grpc

import (
	v1 "github.com/Daniel-Q-Reis/ia_social_media/internal/controller/grpc/v1"
	"github.com/Daniel-Q-Reis/ia_social_media/internal/usecase"
	"github.com/Daniel-Q-Reis/ia_social_media/pkg/logger"
	pbgrpc "google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// NewRouter -.
func NewRouter(app *pbgrpc.Server, t usecase.Translation, l logger.Interface) {
	{
		v1.NewTranslationRoutes(app, t, l)
	}

	reflection.Register(app)
}
