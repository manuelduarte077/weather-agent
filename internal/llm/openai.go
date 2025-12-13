package llm

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"os"
)

type OpenAIResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`

	Error *struct {
		Message string `json:"message"`
		Type    string `json:"type"`
	} `json:"error"`
}

func Ask(prompt string) (string, error) {
	payload := map[string]interface{}{
		"model": "gpt-4.1-mini",
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

	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(
		"POST",
		"https://api.openai.com/v1/chat/completions",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+os.Getenv("OPENAI_API_KEY"))
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result OpenAIResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	// Error de OpenAI
	if result.Error != nil {
		return "", errors.New(result.Error.Message)
	}

	// Sin respuestas
	if len(result.Choices) == 0 {
		return "", errors.New("respuesta vacía del modelo")
	}

	return result.Choices[0].Message.Content, nil
}
