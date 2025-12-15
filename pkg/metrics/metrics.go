package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

// Metrics holds all metrics for the application
type Metrics struct {
	RequestCount    *prometheus.CounterVec
	RequestLatency  *prometheus.HistogramVec
	ErrorCount      *prometheus.CounterVec
	Registry        *prometheus.Registry
	Handler         http.Handler
}

// NewMetrics creates a new Metrics instance
func NewMetrics(namespace string) *Metrics {
	registry := prometheus.NewRegistry()
	
	// Create metrics
	requestCount := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "request_count",
			Help:      "Number of requests received",
		},
		[]string{"method", "endpoint", "status"},
	)
	
	requestLatency := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: namespace,
			Name:      "request_latency_seconds",
			Help:      "Request latency in seconds",
			Buckets:   prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)
	
	errorCount := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "error_count",
			Help:      "Number of errors occurred",
		},
		[]string{"method", "endpoint", "error_type"},
	)
	
	// Register metrics
	registry.MustRegister(requestCount, requestLatency, errorCount)
	
	return &Metrics{
		RequestCount:   requestCount,
		RequestLatency: requestLatency,
		ErrorCount:     errorCount,
		Registry:       registry,
		Handler:        promhttp.HandlerFor(registry, promhttp.HandlerOpts{}),
	}
}

// IncrementRequestCount increments the request count metric
func (m *Metrics) IncrementRequestCount(method, endpoint, status string) {
	m.RequestCount.WithLabelValues(method, endpoint, status).Inc()
}

// ObserveRequestLatency observes the request latency metric
func (m *Metrics) ObserveRequestLatency(method, endpoint string, latency float64) {
	m.RequestLatency.WithLabelValues(method, endpoint).Observe(latency)
}

// IncrementErrorCount increments the error count metric
func (m *Metrics) IncrementErrorCount(method, endpoint, errorType string) {
	m.ErrorCount.WithLabelValues(method, endpoint, errorType).Inc()
}
