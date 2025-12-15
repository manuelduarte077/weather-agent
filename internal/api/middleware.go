package api

import (
	"context"
	"net/http"
	"time"
	"weather-agent/internal/logger"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// requestIDKey is the key used to store request ID in context.
type requestIDKey struct{}

// RequestID returns the request ID from the context.
func RequestID(ctx context.Context) string {
	if id, ok := ctx.Value(requestIDKey{}).(string); ok {
		return id
	}
	return ""
}

// RequestIDMiddleware adds a unique request ID to each request for log correlation.
func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// Add request ID to context
		ctx := context.WithValue(r.Context(), requestIDKey{}, requestID)

		// Add request ID to response header
		w.Header().Set("X-Request-ID", requestID)

		// Continue with the request
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// LoggingMiddleware logs HTTP requests with request ID and timing information.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Wrap response writer to capture status code
		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		// Process request
		next.ServeHTTP(rw, r)

		// Log request details
		duration := time.Since(start)
		requestID := RequestID(r.Context())

		// Use structured logging format
		logHTTPRequest(r.Method, r.URL.Path, rw.statusCode, duration, requestID, r.RemoteAddr)
	})
}

// responseWriter wraps http.ResponseWriter to capture status code.
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// logHTTPRequest logs HTTP request information using structured logging with logrus.
func logHTTPRequest(method, path string, statusCode int, duration time.Duration, requestID, remoteAddr string) {
	fields := logrus.Fields{
		"method":      method,
		"path":        path,
		"status_code": statusCode,
		"status":      http.StatusText(statusCode),
		"duration_ms": duration.Milliseconds(),
		"remote_addr": remoteAddr,
	}

	if requestID != "" {
		fields["request_id"] = requestID
	}

	// Log at appropriate level based on status code
	if statusCode >= 500 {
		// Error level for server errors
		logger.Logger.WithFields(fields).Error("HTTP request")
	} else if statusCode >= 400 {
		// Warning level for client errors
		logger.Logger.WithFields(fields).Warn("HTTP request")
	} else {
		// Info level for successful requests
		logger.Logger.WithFields(fields).Info("HTTP request")
	}
}
