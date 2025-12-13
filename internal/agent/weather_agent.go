package agent

import (
	"encoding/json"
	"fmt"
	"weather-agent/internal/llm"
	"weather-agent/internal/weather"
)

func Run(city string) (map[string]interface{}, error) {
	w, err := weather.GetWeather(city)
	if err != nil {
		return nil, err
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

	response, err := llm.Ask(prompt)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	err = json.Unmarshal([]byte(response), &result)
	if err != nil {
		return nil, err
	}

	result["city"] = city
	result["temp"] = w.Main.Temp
	result["condition"] = w.Weather[0].Description

	return result, nil
}
