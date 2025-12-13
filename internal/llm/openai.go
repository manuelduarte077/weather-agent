package llm

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"os"
)

func Ask(prompt string) (string, error) {
	payload := map[string]interface{}{
		"model": "gpt-4o-mini",
		"messages": []map[string]string{
			{
				"role":    "system",
				"content": "Eres un agente de clima experto en Nicaragua. Responde solo JSON.",
			},
			{
				"role":    "user",
				"content": prompt,
			},
		},
	}

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest(
		"POST",
		"https://api.openai.com/v1/chat/completions",
		bytes.NewBuffer(body),
	)

	req.Header.Set("Authorization", "Bearer "+os.Getenv("OPENAI_API_KEY"))
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	choices := result["choices"].([]interface{})
	if len(choices) == 0 {
		return "", errors.New("sin respuesta del modelo")
	}

	message := choices[0].(map[string]interface{})["message"]
	content := message.(map[string]interface{})["content"]

	return content.(string), nil
}
