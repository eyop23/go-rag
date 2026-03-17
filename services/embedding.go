package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/eyop23/insurance-go/models"
)

func GetEmbedding(text string, apiKeys []string) ([]float64, error) {
	var lastErr error
	for i, apiKey := range apiKeys {
		url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/gemini-embedding-001:embedContent?key=%s", apiKey)

		reqBody := models.GeminiEmbedRequest{}
		reqBody.Content.Parts = []struct {
			Text string `json:"text"`
		}{{Text: text}}

		jsonData, _ := json.Marshal(reqBody)
		resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			lastErr = fmt.Errorf("embedding request failed: %w", err)
			continue
		}

		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		if resp.StatusCode == 429 {
			log.Printf("Embedding: key %d rate limited (429), trying next key", i)
			lastErr = fmt.Errorf("embedding API error %d: %s", resp.StatusCode, string(body))
			continue
		}

		if resp.StatusCode != 200 {
			return nil, fmt.Errorf("embedding API error %d: %s", resp.StatusCode, string(body))
		}

		var result models.GeminiEmbedResponse
		if err := json.Unmarshal(body, &result); err != nil {
			return nil, fmt.Errorf("failed to parse embedding response: %w", err)
		}

		return result.Embedding.Values, nil
	}

	return nil, fmt.Errorf("all embedding API keys exhausted: %w", lastErr)
}
