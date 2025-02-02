package server

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	httpSwagger "github.com/swaggo/http-swagger/v2"

	authHttp "cyansnbrst/auth-service/internal/auth/delivery/http"
	authRepository "cyansnbrst/auth-service/internal/auth/repository"
	authUseCase "cyansnbrst/auth-service/internal/auth/usecase"
)

// Register server handlers
func (s *Server) RegisterHandlers() http.Handler {
	router := httprouter.New()

	// Init repository
	authRepo := authRepository.NewAuthRepository(s.config, s.db)

	// Init use case
	authUC := authUseCase.NewAuthUseCase(s.config, authRepo, s.logger)

	// Init handlers
	authHandlers := authHttp.NewAuthHandlers(s.config, authUC, s.logger)

	// Register auth routes
	authHttp.RegisterAuthRoutes(router, authHandlers)

	// Swagger
	router.ServeFiles("/auth/docs/*filepath", http.Dir("docs"))
	router.HandlerFunc(http.MethodGet, "/auth/swagger/*action", httpSwagger.Handler(
		httpSwagger.URL("/auth/docs/swagger.json"),
	))

	return router
}
