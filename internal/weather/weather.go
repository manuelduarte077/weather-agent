package weather

import "context"

// WeatherService defines the interface for weather data retrieval.
type WeatherService interface {
	GetWeather(ctx context.Context, city string) (*WeatherResponse, error)
}
