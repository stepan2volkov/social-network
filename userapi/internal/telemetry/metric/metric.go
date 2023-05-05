package metric

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	namespace = "userapi"
)

var (
	buckets = []float64{5, 10, 25, 50, 100, 250, 500, 1000, 2500, 5000, 10000}
)

type apiMetrics struct {
	latency *prometheus.HistogramVec
}

func NewAPIMetrics() *apiMetrics {
	latency := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: namespace,
			Name:      "api_latency",
			Help:      "latency of api responses (ms)",
			Buckets:   buckets,
		},
		[]string{"method", "handler", "status"},
	)

	prometheus.MustRegister(latency)

	return &apiMetrics{
		latency: latency,
	}
}

func (am *apiMetrics) WithLatency(method, handler, status string, started time.Time) {
	am.latency.
		With(
			prometheus.Labels{
				"method":  method,
				"handler": handler,
				"status":  status,
			},
		).
		Observe(
			float64(time.Since(started).Milliseconds()),
		)
}
