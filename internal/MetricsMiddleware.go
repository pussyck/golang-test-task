package internal

import (
	"net/http"
	"time"
)

type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (rr *responseRecorder) WriteHeader(code int) {
	rr.statusCode = code
	rr.ResponseWriter.WriteHeader(code)
}

func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rr := &responseRecorder{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(rr, r)

		duration := time.Since(start).Seconds()
		requestsTotal.WithLabelValues(r.Method).Inc()
		requestStatusCodes.WithLabelValues(r.Method, http.StatusText(rr.statusCode)).Inc()
		requestDuration.WithLabelValues(r.Method).Observe(duration)
	})
}
