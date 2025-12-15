package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"
	"weather-agent/internal/agent"
)

// Handler holds dependencies for HTTP handlers.
type Handler struct {
	agent *agent.WeatherAgent
}

// NewHandler creates a new HTTP handler with the provided agent.
func NewHandler(agent *agent.WeatherAgent) *Handler {
	return &Handler{
		agent: agent,
	}
}

// AgentWeatherHandler handles weather agent requests.
// It validates the city parameter, sanitizes input, and processes the request
// with proper context propagation and timeout handling.
func (h *Handler) AgentWeatherHandler(w http.ResponseWriter, r *http.Request) {
	// Create context with timeout for the request
	ctx, cancel := context.WithTimeout(r.Context(), 45*time.Second)
	defer cancel()

	// Only allow GET method
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	city := r.URL.Query().Get("city")
	if city == "" {
		city = "Managua,NI"
	}

	// Validate and sanitize city input
	city, err := ValidateCity(city, 100)
	if err != nil {
		if valErr, ok := err.(*ValidationError); ok {
			http.Error(w, valErr.Message, valErr.StatusCode)
			return
		}
		http.Error(w, "Invalid city parameter", http.StatusBadRequest)
		return
	}

	result, err := h.agent.Run(ctx, city)
	if err != nil {
		logError(r.Context(), "Failed to process weather request: %v", err)
		// Don't expose internal error details to client
		http.Error(w, "Failed to process weather request", http.StatusInternalServerError)
		return
	}

	// Encode JSON to buffer first to avoid writing headers before checking for errors
	jsonData, err := json.Marshal(result)
	if err != nil {
		logError(r.Context(), "Failed to marshal response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	// Set headers only after successful encoding
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusOK)

	// Write the JSON response
	if _, err := w.Write(jsonData); err != nil {
		logError(r.Context(), "Failed to write response: %v", err)
		// Headers already written, can't change status code
		return
	}
}

// logError logs an error with request ID for correlation.
func logError(ctx context.Context, format string, args ...interface{}) {
	requestID := RequestID(ctx)
	log.Printf("[ERROR] [request_id=%s] "+format, append([]interface{}{requestID}, args...)...)
}
