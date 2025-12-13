package api

import (
	"encoding/json"
	"net/http"
	"weather-agent/internal/agent"
)

func AgentWeatherHandler(w http.ResponseWriter, r *http.Request) {
	city := r.URL.Query().Get("city")
	if city == "" {
		city = "Managua,NI"
	}

	result, err := agent.Run(city)
	if err != nil {
		http.Error(w, "Error del agente IA", 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
