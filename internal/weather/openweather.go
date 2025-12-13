package weather

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type WeatherResponse struct {
	Main struct {
		Temp float64 `json:"temp"`
	} `json:"main"`
	Weather []struct {
		Description string `json:"description"`
	} `json:"weather"`
}

func GetWeather(city string) (*WeatherResponse, error) {
	apiKey := os.Getenv("OPENWEATHER_API_KEY")

	url := fmt.Sprintf(
		"https://api.openweathermap.org/data/2.5/weather?q=%s&units=metric&lang=es&appid=%s",
		city,
		apiKey,
	)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data WeatherResponse
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}
