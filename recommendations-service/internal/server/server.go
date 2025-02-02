package server

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"

	"cyansnbrst/recommendations-service/config"
)

// Server struct
type Server struct {
	config      *config.Config
	logger      *zap.Logger
	db          *sql.DB
	redisClient *redis.Client
}

// New server constructor
func NewServer(cfg *config.Config, logger *zap.Logger, db *sql.DB, redisClient *redis.Client) *Server {
	return &Server{
		config:      cfg,
		logger:      logger,
		db:          db,
		redisClient: redisClient,
	}
}

// Run server
func (s *Server) Run() error {
	addr := fmt.Sprintf(":%d", s.config.Port)
	server := &http.Server{
		Addr:         addr,
		Handler:      s.RegisterHandlers(),
		IdleTimeout:  s.config.Timeout.ServerIdle,
		ReadTimeout:  s.config.Timeout.ServerRead,
		WriteTimeout: s.config.Timeout.ServerWrite,
	}

	// Graceful shutdown
	shutDownError := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		sig := <-quit

		s.logger.Info("shutting down server",
			zap.String("signal", sig.String()),
		)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := server.Shutdown(ctx)
		if err != nil {
			shutDownError <- err
		}

		shutDownError <- nil
	}()

	s.logger.Info("starting server",
		zap.String("addr", server.Addr),
		zap.String("env", s.config.Env),
	)

	err := server.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutDownError
	if err != nil {
		return err
	}

	s.logger.Info("stopped server",
		zap.String("addr", server.Addr),
	)

	return nil
}
