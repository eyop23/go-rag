package handlers

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/eyop23/insurance-go/config"
	"github.com/eyop23/insurance-go/models"
	"github.com/eyop23/insurance-go/services"
)

func AskHandler(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.AskRequest
		if err := c.ShouldBindJSON(&req); err != nil || req.Query == "" {
			c.JSON(http.StatusBadRequest, models.AskResponse{Error: "Query is required"})
			return
		}

		log.Printf("Query: %q", req.Query)

		// Step 1: Generate embedding
		log.Println("Step 1: Generating query embedding...")
		embedding, err := services.GetEmbedding(req.Query, cfg.GoogleAPIKeys)
		if err != nil {
			log.Printf("Embedding error: %v", err)
			if strings.Contains(err.Error(), "keys exhausted") || strings.Contains(err.Error(), "429") {
				c.JSON(http.StatusTooManyRequests, models.AskResponse{Error: "Gemini API rate limit reached. Please update your API key."})
			} else {
				c.JSON(http.StatusServiceUnavailable, models.AskResponse{Error: "Service is temporarily busy. Please try again in a moment."})
			}
			return
		}
		log.Printf("Embedding done: %d dimensions", len(embedding))

		// Step 2: Query Pinecone
		log.Println("Step 2: Querying Pinecone...")
		matches, err := services.QueryPinecone(embedding, cfg.PineconeHost, cfg.PineconeKey)
		if err != nil {
			log.Printf("Pinecone error: %v", err)
			c.JSON(http.StatusServiceUnavailable, models.AskResponse{Answer: "Service is temporarily busy. Please try again in a moment."})
			return
		}
		log.Printf("Found %d matches", len(matches))

		if len(matches) == 0 {
			c.JSON(http.StatusOK, models.AskResponse{Answer: "No relevant information found."})
			return
		}

		// Step 3: Build context
		log.Println("Step 3: Building context...")
		var contextParts []string
		for _, m := range matches {
			if text, ok := m.Metadata["text"]; ok {
				contextParts = append(contextParts, text)
			}
		}
		context := strings.Join(contextParts, "\n\n")

		// Step 4: Get answer from LLM
		log.Println("Step 4: Generating answer...")
		answer, err := services.GetAnswerFromLLM(req.Query, context, cfg.GeminiAPIURL, cfg.GoogleAPIKeys)
		if err != nil {
			log.Printf("LLM error: %v", err)
			if strings.Contains(err.Error(), "keys exhausted") || strings.Contains(err.Error(), "429") {
				c.JSON(http.StatusTooManyRequests, models.AskResponse{Error: "Gemini API rate limit reached. Please update your API key."})
			} else {
				c.JSON(http.StatusServiceUnavailable, models.AskResponse{Error: "Service is temporarily busy. Please try again in a moment."})
			}
			return
		}

		log.Printf("Answer: %.100s...", answer)

		c.JSON(http.StatusOK, models.AskResponse{Answer: answer})
	}
}
