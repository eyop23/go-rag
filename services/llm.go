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

func GetAnswerFromLLM(query, context, apiURL string, apiKeys []string) (string, error) {
	prompt := fmt.Sprintf(`You are a knowledgeable book recommendation assistant. Answer the question based ONLY on the following context.
If the question asks to list, filter, or compare multiple books, scan ALL the records in the context and include every matching result.
Do not skip any matching record. If no records match, say so clearly.

Context:
%s

Question:
%s`, context, query)

	reqBody := models.GeminiGenRequest{}
	reqBody.Contents = []struct {
		Parts []struct {
			Text string `json:"text"`
		} `json:"parts"`
	}{
		{
			Parts: []struct {
				Text string `json:"text"`
			}{{Text: prompt}},
		},
	}

	jsonData, _ := json.Marshal(reqBody)

	var lastErr error
	for i, apiKey := range apiKeys {
		url := fmt.Sprintf("%s?key=%s", apiURL, apiKey)

		resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			lastErr = fmt.Errorf("LLM request failed: %w", err)
			continue
		}

		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		if resp.StatusCode == 429 {
			log.Printf("LLM: key %d rate limited (429), trying next key", i)
			lastErr = fmt.Errorf("gemini LLM error %d: %s", resp.StatusCode, string(body))
			continue
		}

		if resp.StatusCode != 200 {
			return "", fmt.Errorf("gemini LLM error %d: %s", resp.StatusCode, string(body))
		}

		var result models.GeminiGenResponse
		json.Unmarshal(body, &result)

		if len(result.Candidates) > 0 && len(result.Candidates[0].Content.Parts) > 0 {
			return result.Candidates[0].Content.Parts[0].Text, nil
		}
		return "No response from Gemini", nil
	}

	return "", fmt.Errorf("all LLM API keys exhausted: %w", lastErr)
}
