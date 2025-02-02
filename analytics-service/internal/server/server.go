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

	"go.uber.org/zap"

	"cyansnbrst/analytics-service/config"
)

// Server struct
type Server struct {
	config *config.Config
	logger *zap.Logger
	db     *sql.DB
}

// New server constructor
func NewServer(cfg *config.Config, logger *zap.Logger, db *sql.DB) *Server {
	return &Server{
		config: cfg,
		logger: logger,
		db:     db,
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

		ctx, cancel := context.WithTimeout(context.Background(), s.config.Timeout.ServerShutdown)
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
