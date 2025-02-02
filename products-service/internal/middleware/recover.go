package middleware

import (
	"fmt"
	"net/http"

	"go.uber.org/zap"

	erp "cyansnbrst/products-service/pkg/error_responses"
)

// Panic recoverer middleware
func (mw *MiddlewareManager) RecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				mw.logger.Error("recovered from panic", zap.Any("error", err))
				erp.ServerErrorResponse(w, r, mw.logger, fmt.Errorf("%s", err))
				return
			}
		}()

		next.ServeHTTP(w, r)
	})
}
