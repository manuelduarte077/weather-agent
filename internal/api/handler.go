package api

import (
	"encoding/json"
	"log"
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
		log.Println("AGENT ERROR:", err.Error())

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
