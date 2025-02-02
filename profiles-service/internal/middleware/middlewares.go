package middleware

import (
	"go.uber.org/zap"

	"cyansnbrst/profiles-service/config"
)

// Middleware manager struct
type MiddlewareManager struct {
	cfg    *config.Config
	logger *zap.Logger
}

// Middleware manager constructor
func NewMiddlewareManager(cfg *config.Config, logger *zap.Logger) *MiddlewareManager {
	return &MiddlewareManager{cfg: cfg, logger: logger}
}
