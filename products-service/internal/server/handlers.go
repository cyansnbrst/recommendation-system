package server

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	httpSwagger "github.com/swaggo/http-swagger/v2"

	"cyansnbrst/products-service/internal/middleware"
	productsHttp "cyansnbrst/products-service/internal/products/delivery/http"
	productsRepository "cyansnbrst/products-service/internal/products/repository"
	productsUseCase "cyansnbrst/products-service/internal/products/usecase"
)

// Register server handlers
func (s *Server) RegisterHandlers() http.Handler {
	router := httprouter.New()

	// Init repository
	productsRepo := productsRepository.NewProductsRepository(s.config, s.db)

	// Init use case
	productsUC := productsUseCase.NewProductsUseCase(s.config, productsRepo, s.logger)

	// Init handlers
	productsHandlers := productsHttp.NewProductsHandlers(s.config, productsUC, s.logger, s.kafkaUserWriter, s.kafkaProductWriter)

	// Init middleware
	mw := middleware.NewMiddlewareManager(s.config, s.logger)

	// Register products routes
	productsHttp.RegisterProductsRoutes(router, productsHandlers, mw)

	// Swagger
	router.ServeFiles("/products/docs/*filepath", http.Dir("docs"))
	router.HandlerFunc(http.MethodGet, "/products/swagger/*action", httpSwagger.Handler(
		httpSwagger.URL("/products/docs/swagger.json"),
	))

	return mw.RecoverPanic(mw.Authenticate(router))
}
