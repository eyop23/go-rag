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
	prompt := fmt.Sprintf(`You are a book recommendation assistant. You ONLY answer questions about books. For greetings and small talk, respond warmly and briefly.

If the user asks about anything unrelated to books (people, sports, music, movies, history, etc.), respond with something like: "Hey, I'm Go-RAG, a book recommender! I can only help with book-related questions. Are you looking for a book recommendation?"

For book-related questions, answer based on the following context. If the question asks to list, filter, or compare multiple books, scan ALL the records and include every matching result. If no records match, say so clearly.

Context:
%s

User message:
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
