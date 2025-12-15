package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"weather-agent/internal/agent"
	"weather-agent/internal/api"
	"weather-agent/internal/config"
	"weather-agent/internal/llm"
	"weather-agent/internal/logger"
	"weather-agent/internal/weather"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logger.Logger.WithError(err).Fatal("Failed to load configuration")
	}

	// Initialize services with dependency injection
	weatherClient := weather.NewOpenWeatherClient(cfg.OpenWeatherAPIKey)
	llmClient := llm.NewOpenAIClient(cfg.OpenAIAPIKey)
	weatherAgent := agent.NewWeatherAgent(weatherClient, llmClient)

	// Create HTTP handler
	handler := api.NewHandler(weatherAgent)

	// Setup HTTP routes with middleware
	mux := http.NewServeMux()
	mux.HandleFunc("/agent/weather", handler.AgentWeatherHandler)

	// Add health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Wrap mux with middleware for request ID and logging
	handlerWithMiddleware := api.RequestIDMiddleware(api.LoggingMiddleware(mux))

	// Create HTTP server with timeouts
	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      handlerWithMiddleware,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		logger.Logger.WithField("port", cfg.Port).Info("Weather Agent running")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Logger.WithError(err).Fatal("Server failed to start")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Logger.Info("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Logger.WithError(err).Fatal("Server forced to shutdown")
	}

	logger.Logger.Info("Server exited gracefully")

}
