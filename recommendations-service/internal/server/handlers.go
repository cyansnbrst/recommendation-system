package server

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"

	"cyansnbrst/recommendations-service/internal/client"
	"cyansnbrst/recommendations-service/internal/middleware"
	"cyansnbrst/recommendations-service/internal/recommendations/delivery/consumers"
	recommendationsHttp "cyansnbrst/recommendations-service/internal/recommendations/delivery/http"
	recommendationsRepository "cyansnbrst/recommendations-service/internal/recommendations/repository"
	recommendationsUseCase "cyansnbrst/recommendations-service/internal/recommendations/usecase"
	"cyansnbrst/recommendations-service/pkg/metric"
)

// Register server handlers
func (s *Server) RegisterHandlers() http.Handler {
	metrics, err := metric.CreateMetrics(s.config.Metrics.URL, s.config.Metrics.ServiceName)
	if err != nil {
		s.logger.Error("create metrics error", zap.Error(err))
	}

	router := httprouter.New()

	// Init repository
	recommendationsRepo := recommendationsRepository.NewRecommendationsRepository(s.config, s.db)
	recommendationsRedisRepo := recommendationsRepository.NewRecommendationsRedisRepository(s.config, s.redisClient)

	// Init use case
	recommendationsUC := recommendationsUseCase.NewRecommendationsUseCase(s.config, recommendationsRepo, recommendationsRedisRepo, s.logger)

	// Init handlers
	recommendationsHandlers := recommendationsHttp.NewRecommendationsHandlers(s.config, recommendationsUC, s.logger)

	// Init middleware
	mw := middleware.NewMiddlewareManager(s.config, s.logger)

	// Register recommendations routes
	recommendationsHttp.RegisterRecommendationsRoutes(router, recommendationsHandlers, mw)

	// Init kafka consumers
	kafkaClient := client.NewKafkaClient(s.config, s.logger)
	kafkaHandlers := consumers.NewKafkaMessageHandlers(s.config, recommendationsUC, s.logger)

	kafkaClient.AddReader("product", s.config.Kafka.GroupID, kafkaHandlers.HandleProductMessage)
	kafkaClient.AddReader("user", s.config.Kafka.GroupID, kafkaHandlers.HandleUserMessage)
	kafkaClient.Run()

	// Swagger
	router.ServeFiles("/recommendations/docs/*filepath", http.Dir("docs"))
	router.HandlerFunc(http.MethodGet, "/recommendations/swagger/*action", httpSwagger.Handler(
		httpSwagger.URL("/recommendations/docs/swagger.json"),
	))

	wrappedRouter := mw.MetricsMiddleware(metrics)(router)

	return mw.RecoverPanic(mw.Authenticate(wrappedRouter))
}
