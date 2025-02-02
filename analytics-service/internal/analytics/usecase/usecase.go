package usecase

import (
	"time"

	"go.uber.org/zap"

	"cyansnbrst/analytics-service/config"
	"cyansnbrst/analytics-service/internal/analytics"
)

// Analytics UseCase struct
type analyticsUC struct {
	cfg           *config.Config
	analyticsRepo analytics.Repository
	logger        *zap.Logger
}

// New analytics constructor
func NewAnalyticsUseCase(cfg *config.Config, analyticsRepo analytics.Repository, logger *zap.Logger) analytics.UseCase {
	return &analyticsUC{cfg: cfg, analyticsRepo: analyticsRepo, logger: logger}
}

// Generate recommendations for user
func (u *analyticsUC) Insert(action string, objectID string, actionTime time.Time) error {
	return u.analyticsRepo.Insert(action, objectID, actionTime)
}
