package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"social/api/config"
	v1 "social/api/internal/controller/http/v1"
	"social/api/internal/repo/postgres"
	"social/api/internal/usecase"
)

func main() {
	// Initialize logger
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).With().Timestamp().Logger()

	// Load configuration
	cfg := config.MustLoad()

	// Connect to database
	pool, err := pgxpool.New(context.Background(), cfg.URL)
	if err != nil {
		logger.Fatal().Err(err).Msg("Unable to connect to database")
	}
	defer pool.Close()

	// Initialize repositories
	userRepo := postgres.NewUserRepo(pool)
	postRepo := postgres.NewPostRepo(pool)
	commentRepo := postgres.NewCommentRepo(pool)
	likeRepo := postgres.NewLikeRepo(pool)
	followRepo := postgres.NewFollowRepo(pool)

	// Initialize use cases
	userUseCase := usecase.NewUserUseCase(userRepo)
	postUseCase := usecase.NewPostUseCase(postRepo, userRepo)
	commentUseCase := usecase.NewCommentUseCase(commentRepo, userRepo, postRepo)
	interactionUseCase := usecase.NewInteractionUseCase(likeRepo, followRepo, userRepo)

	// Initialize handler
	handler := v1.NewHandler(userUseCase, postUseCase, commentUseCase, interactionUseCase)

	// Initialize router
	r := chi.NewRouter()

	// Add middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(middleware.Timeout(60 * time.Second))

	// Register routes
	handler.RegisterRoutes(r)

	// Start server
	server := &http.Server{
		Addr:         cfg.Address,
		Handler:      r,
		ReadTimeout:  time.Duration(cfg.Timeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Timeout) * time.Second,
		IdleTimeout:  time.Duration(cfg.IdleTimeout) * time.Second,
	}

	// Server run context
	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	// Listen for syscall signals for process to interrupt/quit
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig

		// Trigger graceful shutdown
		shutdownCtx, cancel := context.WithTimeout(serverCtx, 30*time.Second)
		defer cancel()

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				logger.Fatal().Msg("graceful shutdown timed out.. forcing exit.")
			}
		}()

		// Trigger graceful shutdown
		err := server.Shutdown(shutdownCtx)
		if err != nil {
			logger.Fatal().Err(err).Msg("server shutdown failed")
		}

		// Cancel server context to close database connections
		serverStopCtx()
	}()

	// Run the server
	logger.Info().Str("address", cfg.Address).Msg("server started")
	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		logger.Fatal().Err(err).Msg("server startup failed")
	}

	// Wait for server context to be stopped
	<-serverCtx.Done()

	logger.Info().Msg("server exited properly")
}