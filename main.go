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

	http.HandleFunc("/agent/weather", api.AgentWeatherHandler)

	log.Println("🤖 Weather Agent corriendo en puerto", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))

}
