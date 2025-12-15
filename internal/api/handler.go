package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
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

	// Sanitize and validate city input
	city = strings.TrimSpace(city)
	if city == "" {
		http.Error(w, "Invalid city parameter", http.StatusBadRequest)
		return
	}

	// Limit city length to prevent abuse
	if len(city) > 100 {
		http.Error(w, "City parameter too long", http.StatusBadRequest)
		return
	}

	result, err := h.agent.Run(ctx, city)
	if err != nil {
		requestID := RequestID(r.Context())
		log.Printf("[ERROR] [request_id=%s] Failed to process weather request: %v", requestID, err)

		// Don't expose internal error details to client
		http.Error(w, "Failed to process weather request", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	if err := json.NewEncoder(w).Encode(result); err != nil {
		requestID := RequestID(r.Context())
		log.Printf("[ERROR] [request_id=%s] Failed to encode response: %v", requestID, err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
