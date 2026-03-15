package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/eyop23/insurance-go/models"
)

func GetAnswerFromLLM(query, context, apiURL, apiKey string) (string, error) {
	url := fmt.Sprintf("%s?key=%s", apiURL, apiKey)

	prompt := fmt.Sprintf(`You are a knowledgeable Ethiopian music assistant. Answer the question based ONLY on the following context.
If the question asks to list, filter, or compare multiple artists, scan ALL the records in the context and include every matching result.
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
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("LLM request failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
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
