package weather

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"
)

// WeatherResponse represents the weather data structure from OpenWeatherMap API.
type WeatherResponse struct {
	Main struct {
		Temp float64 `json:"temp"`
	} `json:"main"`
	Weather []struct {
		Description string `json:"description"`
	} `json:"weather"`
}

// OpenWeatherClient implements WeatherService using OpenWeatherMap API.
type OpenWeatherClient struct {
	apiKey     string
	httpClient *http.Client
	baseURL    string
}

// NewOpenWeatherClient creates a new OpenWeatherMap client with proper configuration.
func NewOpenWeatherClient(apiKey string) *OpenWeatherClient {
	if apiKey == "" {
		apiKey = os.Getenv("OPENWEATHER_API_KEY")
	}
	return &OpenWeatherClient{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		baseURL: "https://api.openweathermap.org/data/2.5/weather",
	}
}

// GetWeather retrieves weather information for the specified city.
// It validates the API key and HTTP response status, returning wrapped errors for better traceability.
func (c *OpenWeatherClient) GetWeather(ctx context.Context, city string) (*WeatherResponse, error) {
	if c.apiKey == "" {
		return nil, fmt.Errorf("openweathermap API key is not set")
	}

	if city == "" {
		return nil, fmt.Errorf("city parameter cannot be empty")
	}

	// Sanitize city parameter by URL encoding
	encodedCity := url.QueryEscape(city)

	reqURL := fmt.Sprintf(
		"%s?q=%s&units=metric&lang=es&appid=%s",
		c.baseURL,
		encodedCity,
		c.apiKey,
	)

	req, err := http.NewRequestWithContext(ctx, "GET", reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch weather data: %w", err)
	}
	defer resp.Body.Close()

	// Validate HTTP status code
	if resp.StatusCode != http.StatusOK {
		var errorResp struct {
			Message string `json:"message"`
			Code    int    `json:"cod"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&errorResp); err == nil {
			return nil, fmt.Errorf("openweathermap API error (code %d): %s", errorResp.Code, errorResp.Message)
		}
		return nil, fmt.Errorf("openweathermap API returned status %d", resp.StatusCode)
	}

	var data WeatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode weather response: %w", err)
	}

	// Validate that weather data is present
	if len(data.Weather) == 0 {
		return nil, fmt.Errorf("weather data is incomplete: no weather conditions found")
	}

	return &data, nil
}
