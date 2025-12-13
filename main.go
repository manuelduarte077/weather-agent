package main

import (
	"log"
	"net/http"
	"os"
	"weather-agent/internal/api"
)

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/agent/weather", func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if recover() != nil {
				http.Error(w, "Error interno del agente", 500)
			}
		}()
		api.AgentWeatherHandler(w, r)
	})

	log.Println("🤖 Weather Agent corriendo en puerto", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))

}
