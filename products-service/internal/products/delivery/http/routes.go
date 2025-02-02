package http

import (
	"net/http"

	"github.com/julienschmidt/httprouter"

	"cyansnbrst/products-service/internal/middleware"
	"cyansnbrst/products-service/internal/products"
)

// Register products routes
func RegisterProductsRoutes(router *httprouter.Router, h products.Handlers, mw *middleware.MiddlewareManager) {
	router.HandlerFunc(http.MethodPost, "/products/create", mw.RequireAdminRights(h.Create()))
	router.HandlerFunc(http.MethodDelete, "/products/delete/:id", mw.RequireAdminRights(h.Delete()))
	router.HandlerFunc(http.MethodPut, "/products/update/:id", mw.RequireAdminRights(h.Update()))
	router.HandlerFunc(http.MethodGet, "/products/view/:id", mw.RequireAuthenticatedUser(h.Get()))
}
