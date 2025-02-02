package consumers

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"

	"cyansnbrst/recommendations-service/config"
	"cyansnbrst/recommendations-service/internal/models"
	"cyansnbrst/recommendations-service/internal/recommendations"
)

// Kafka message handlers struct
type KafkaMessageHandlers struct {
	cfg               *config.Config
	recommendationsUC recommendations.UseCase
	logger            *zap.Logger
}

// Kafka message handlers constructor
func NewKafkaMessageHandlers(cfg *config.Config, recommendationsUC recommendations.UseCase, logger *zap.Logger) *KafkaMessageHandlers {
	return &KafkaMessageHandlers{
		cfg:               cfg,
		recommendationsUC: recommendationsUC,
		logger:            logger,
	}
}

// Kafka product message handler
func (h *KafkaMessageHandlers) HandleProductMessage(msg kafka.Message) error {
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

	productID, _ := strconv.Atoi(string(msg.Key))
	tags := payload.Tags

	switch payload.Action {
	case "view_products":
		err := h.recommendationsUC.IncrementPopularity(int64(productID))
		if err != nil {
			h.logger.Error("failed to increment popularity", zap.Error(err))
			return err
		}
	case "product_create":
		err := h.recommendationsUC.InsertProduct(int64(productID), tags)
		if err != nil {
			h.logger.Error("failed to create a product", zap.Error(err))
			return err
		}
		err = h.recommendationsUC.UpdateRecommendationsForProduct(int64(productID), tags)
		if err != nil {
			h.logger.Error("failed to generate product recommendations", zap.Error(err))
			return err
		}
	case "product_update":
		err = h.recommendationsUC.UpdateRecommendationsForProduct(int64(productID), tags)
		if err != nil {
			h.logger.Error("failed to generate product recommendations", zap.Error(err))
			return err
		}
	case "product_delete":
		err = h.recommendationsUC.DeleteProduct(int64(productID))
		if err != nil {
			h.logger.Error("failed to delete product", zap.Error(err))
			return err
		}
	default:
		h.logger.Warn("unrecognized action", zap.String("action", payload.Action))
	}

	return nil
}

// Kafka product message handler
func (h *KafkaMessageHandlers) HandleUserMessage(msg kafka.Message) error {
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

	userUID := msg.Key
	tags := payload.Tags

	switch payload.Action {
	case "user_update":
		err := h.recommendationsUC.GenerateRecommendationsForUser(string(userUID), tags)
		if err != nil {
			h.logger.Error("failed to generate recommendations", zap.Error(err))
			return err
		}
	default:
		h.logger.Warn("unrecognized action", zap.String("action", payload.Action))
	}

	return nil
}
