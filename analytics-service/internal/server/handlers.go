package server

import (
	"net/http"

	"github.com/julienschmidt/httprouter"

	"cyansnbrst/analytics-service/internal/analytics/delivery/consumers"
	analyticsRepository "cyansnbrst/analytics-service/internal/analytics/repository"
	analyticsUseCase "cyansnbrst/analytics-service/internal/analytics/usecase"
	"cyansnbrst/analytics-service/internal/client"
)

// Register server handlers
func (s *Server) RegisterHandlers() http.Handler {
	router := httprouter.New()

	// Init repository
	analyticsRepo := analyticsRepository.NewAnalyticsRepository(s.config, s.db)

	// Init use case
	analyticsUC := analyticsUseCase.NewAnalyticsUseCase(s.config, analyticsRepo, s.logger)

	// Init kafka consumers
	kafkaClient := client.NewKafkaClient(s.config, s.logger)
	kafkaHandlers := consumers.NewKafkaMessageHandlers(s.config, analyticsUC, s.logger)

	kafkaClient.AddReader("product", s.config.Kafka.GroupID, kafkaHandlers.HandleMessage)
	kafkaClient.AddReader("user", s.config.Kafka.GroupID, kafkaHandlers.HandleMessage)
	kafkaClient.Run()

	return router
}
