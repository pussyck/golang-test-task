package internal

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

var (
	requestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of requests for each HTTP method",
		},
		[]string{"method"},
	)

	requestStatusCodes = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_status_codes_total",
			Help: "Total number of requests by HTTP status codes",
		},
		[]string{"method", "status_code"},
	)

	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)
)

func InitMetrics() {
	prometheus.MustRegister(requestsTotal)
	prometheus.MustRegister(requestStatusCodes)
	prometheus.MustRegister(requestDuration)
}

func HandleMetrics() http.Handler {
	return promhttp.Handler()
}
