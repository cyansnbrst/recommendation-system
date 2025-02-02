package middleware

import (
	"net/http"
	"time"

	"cyansnbrst/recommendations-service/pkg/metric"
)

// Prometheus metrics middleware
func (mw *MiddlewareManager) MetricsMiddleware(metrics metric.Metrics) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			next.ServeHTTP(w, r)

			metrics.ObserveResponseTime(r.Method, r.URL.Path, time.Since(start).Seconds())
			metrics.IncHits(r.Method, r.URL.Path)
		})
	}
}
