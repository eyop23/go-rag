package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/eyop23/insurance-go/models"
)

func QueryPinecone(embedding []float64, host, apiKey string) ([]models.PineconeMatch, error) {
	url := host + "/query"

	reqBody := models.PineconeQueryRequest{
		Vector:          embedding,
		TopK:            20,
		IncludeMetadata: true,
	}

	jsonData, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	req.Header.Set("Api-Key", apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("pinecone query failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("pinecone query error %d: %s", resp.StatusCode, string(body))
	}

	var result models.PineconeQueryResponse
	json.Unmarshal(body, &result)
	return result.Matches, nil
}

func GetPineconeVectorCount(host, apiKey string) (int, error) {
	url := host + "/describe_index_stats"

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer([]byte("{}")))
	req.Header.Set("Api-Key", apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("pinecone stats failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return 0, fmt.Errorf("pinecone stats error %d: %s", resp.StatusCode, string(body))
	}

	var result struct {
		TotalVectorCount int `json:"totalVectorCount"`
	}
	json.Unmarshal(body, &result)
	return result.TotalVectorCount, nil
}

func UpsertToPinecone(vectors []models.PineconeVector, host, apiKey string) error {
	url := host + "/vectors/upsert"

	reqBody := models.PineconeUpsertRequest{Vectors: vectors}
	jsonData, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	req.Header.Set("Api-Key", apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("pinecone upsert failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return fmt.Errorf("pinecone error %d: %s", resp.StatusCode, string(body))
	}

	return nil
}
