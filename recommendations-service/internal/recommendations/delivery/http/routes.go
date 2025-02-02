package http

import (
	"net/http"

	"github.com/julienschmidt/httprouter"

	"cyansnbrst/recommendations-service/internal/middleware"
	"cyansnbrst/recommendations-service/internal/recommendations"
)

// Register recommendations routes
func RegisterRecommendationsRoutes(router *httprouter.Router, h recommendations.Handlers, mw *middleware.MiddlewareManager) {
	router.HandlerFunc(http.MethodGet, "/recommendations", mw.RequireAuthenticatedUser(h.GetInfo()))
}
