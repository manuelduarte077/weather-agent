package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"weather-agent/internal/llm"
	"weather-agent/internal/weather"
)

// WeatherAgent handles weather analysis using weather and LLM services.
type WeatherAgent struct {
	weatherService weather.WeatherService
	llmService     llm.LLMService
}

// NewWeatherAgent creates a new WeatherAgent with the provided services.
func NewWeatherAgent(weatherService weather.WeatherService, llmService llm.LLMService) *WeatherAgent {
	return &WeatherAgent{
		weatherService: weatherService,
		llmService:     llmService,
	}
}

// Run analyzes the weather for the given city and returns recommendations.
// It uses the weather service to get current conditions and the LLM service
// to generate analysis and recommendations.
func (a *WeatherAgent) Run(ctx context.Context, city string) (map[string]interface{}, error) {
	if city == "" {
		return nil, fmt.Errorf("city parameter cannot be empty")
	}

	w, err := a.weatherService.GetWeather(ctx, city)
	if err != nil {
		return nil, fmt.Errorf("failed to get weather data: %w", err)
	}

	if len(w.Weather) == 0 {
		return nil, fmt.Errorf("weather data incomplete: no weather conditions available")
	}

	prompt := fmt.Sprintf(`
		Contexto:
		País: Nicaragua
		Ciudad: %s
		Temperatura: %.1f °C
		Condición: %s

		Instrucciones:
		- Analiza el clima
		- Da una recomendación práctica
		- Evalúa riesgo

		Responde SOLO este JSON:
		{
			"analysis": "...",
			"recommendation": "...",
			"risk_level": "bajo | medio | alto"
		}`,
		city,
		w.Main.Temp,
		w.Weather[0].Description,
	)

	response, err := a.llmService.Ask(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("failed to get LLM response: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal([]byte(response), &result); err != nil {
		return nil, fmt.Errorf("failed to parse LLM JSON response: %w", err)
	}

	result["city"] = city
	result["temp"] = w.Main.Temp
	result["condition"] = w.Weather[0].Description

	return result, nil
}
