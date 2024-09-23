package internal

import (
	"net/http"
	"time"
)

// responseRecorder записывает статус код и размер ответа
type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (rr *responseRecorder) WriteHeader(code int) {
	rr.statusCode = code
	rr.ResponseWriter.WriteHeader(code)
}

// MetricsMiddleware обновляет метрики для каждого запроса
func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rr := &responseRecorder{ResponseWriter: w, statusCode: http.StatusOK}

		path := r.URL.Path

		next.ServeHTTP(rr, r)

		duration := time.Since(start).Seconds()
		method := r.Method
		statusCode := rr.statusCode

		requestsTotal.WithLabelValues(method, path).Inc()
		requestStatusCodes.WithLabelValues(method, path, http.StatusText(statusCode)).Inc()
		requestDuration.WithLabelValues(method, path).Observe(duration)
	})
}
