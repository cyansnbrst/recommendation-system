package http

import (
	"net/http"

	"github.com/julienschmidt/httprouter"

	"cyansnbrst/auth-service/internal/auth"
)

// Register auth routes
func RegisterAuthRoutes(router *httprouter.Router, h auth.Handlers) {
	router.HandlerFunc(http.MethodPost, "/auth/login", h.Login())
	router.HandlerFunc(http.MethodPost, "/auth/register", h.Register())
	router.HandlerFunc(http.MethodGet, "/auth/authenticate", h.TokenValidation())
}
