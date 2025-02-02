package server

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	httpSwagger "github.com/swaggo/http-swagger/v2"

	"cyansnbrst/profiles-service/internal/middleware"
	profilesHttp "cyansnbrst/profiles-service/internal/profiles/delivery/http"
	profilesRepository "cyansnbrst/profiles-service/internal/profiles/repository"
	profilesUseCase "cyansnbrst/profiles-service/internal/profiles/usecase"
)

// Register server handlers
func (s *Server) RegisterHandlers() http.Handler {
	router := httprouter.New()

	// Init repository
	profilesRepo := profilesRepository.NewProfilesRepository(s.config, s.db)

	// Init use case
	profilesUC := profilesUseCase.NewProfilesUseCase(s.config, profilesRepo, s.logger)

	// Init handlers
	profilesHandlers := profilesHttp.NewProfilesHandlers(s.config, profilesUC, s.logger, s.kafkaWriter)

	// Init middleware
	mw := middleware.NewMiddlewareManager(s.config, s.logger)

	// Register profiles routes
	profilesHttp.RegisterProfileRoutes(router, profilesHandlers, mw)

	// Swagger
	router.ServeFiles("/profiles/docs/*filepath", http.Dir("docs"))
	router.HandlerFunc(http.MethodGet, "/profiles/swagger/*action", httpSwagger.Handler(
		httpSwagger.URL("/profiles/docs/swagger.json"),
	))

	return mw.RecoverPanic(mw.Authenticate(router))
}
