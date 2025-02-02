package consumers

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"

	"cyansnbrst/analytics-service/config"
	"cyansnbrst/analytics-service/internal/analytics"
	"cyansnbrst/analytics-service/internal/models"
)

// Kafka message handlers struct
type KafkaMessageHandlers struct {
	cfg         *config.Config
	analyticsUC analytics.UseCase
	logger      *zap.Logger
}

// Kafka message handlers constructor
func NewKafkaMessageHandlers(cfg *config.Config, analyticsUC analytics.UseCase, logger *zap.Logger) analytics.KafkaHandlers {
	return &KafkaMessageHandlers{
		cfg:         cfg,
		analyticsUC: analyticsUC,
		logger:      logger,
	}
}

// Kafka message handler
func (h *KafkaMessageHandlers) HandleMessage(msg kafka.Message) error {
	var payload models.KafkaMessageDTO

	err := json.Unmarshal(msg.Value, &payload)
	if err != nil {
		h.logger.Error("failed to unmarshal Kafka message", zap.Error(err))
		return err
	}

	h.logger.Info("kafka message received",
		zap.ByteString("key", msg.Key),
		zap.String("action", payload.Action),
		zap.Strings("tags", payload.Tags),
		zap.String("time", payload.Time),
	)
	payload.Action = strings.TrimSpace(strings.ToLower(strings.TrimRight(payload.Action, "\n\r")))

	objectID := string(msg.Key)
	action := payload.Action
	actionTime := payload.Time
	parsedTime, err := time.Parse(time.RFC3339, actionTime)
	if err != nil {
		h.logger.Error("couldn't parse time", zap.Error(err))
		return err
	}

	err = h.analyticsUC.Insert(action, objectID, parsedTime)
	if err != nil {
		h.logger.Error("failed to insert action", zap.Error(err))
		return err
	}

	return nil
}
