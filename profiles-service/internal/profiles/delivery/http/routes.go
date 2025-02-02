package http

import (
	"net/http"

	"github.com/julienschmidt/httprouter"

	"cyansnbrst/profiles-service/internal/middleware"
	"cyansnbrst/profiles-service/internal/profiles"
)

// Register profile routes
func RegisterProfileRoutes(router *httprouter.Router, h profiles.Handlers, mw *middleware.MiddlewareManager) {
	router.HandlerFunc(http.MethodGet, "/profiles", mw.RequireAuthenticatedUser(h.GetInfo()))
	router.HandlerFunc(http.MethodPut, "/profiles/edit", mw.RequireAuthenticatedUser(h.EditData()))
	router.HandlerFunc(http.MethodPost, "/profiles/create/:uid", h.CreateProfile())
}
