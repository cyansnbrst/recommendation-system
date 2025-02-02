package metric

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// App metrics interface
type Metrics interface {
	IncHits(method, path string)
	ObserveResponseTime(method, path string, observeTime float64)
}

// Prometheus metrics struct
type PrometheusMetrics struct {
	HitsTotal prometheus.Counter
	Hits      *prometheus.CounterVec
	Times     *prometheus.HistogramVec
}

// Create metrics with address and name
func CreateMetrics(address string, name string) (Metrics, error) {
	var metr PrometheusMetrics
	metr.HitsTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: name + "_hits_total",
	})

	if err := prometheus.Register(metr.HitsTotal); err != nil {
		return nil, err
	}

	metr.Hits = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: name + "_hits",
		},
		[]string{"method", "path"},
	)

	if err := prometheus.Register(metr.Hits); err != nil {
		return nil, err
	}

	metr.Times = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: name + "_times",
		},
		[]string{"method", "path"},
	)

	if err := prometheus.Register(metr.Times); err != nil {
		return nil, err
	}

	if err := prometheus.Register(collectors.NewBuildInfoCollector()); err != nil {
		return nil, err
	}

	go func() {
		router := httprouter.New()
		router.GET("/metrics", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
			promhttp.Handler().ServeHTTP(w, r)
		})
		log.Printf("metrics server is running on port: %s", address)
		if err := http.ListenAndServe(address, router); err != nil {
			log.Fatal(err)
		}
	}()

	return &metr, nil
}

// IncHits
func (metr *PrometheusMetrics) IncHits(method, path string) {
	metr.HitsTotal.Inc()
	metr.Hits.WithLabelValues(method, path).Inc()
}

// ObserveResponseTime
func (metr *PrometheusMetrics) ObserveResponseTime(method, path string, observeTime float64) {
	metr.Times.WithLabelValues(method, path).Observe(observeTime)
}
