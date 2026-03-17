package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	GoogleAPIKeys []string
	GeminiAPIURL  string
	PineconeKey   string
	PineconeHost  string
	Port          string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env not found, using system env")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8090"
	}

	var apiKeys []string
	if multi := os.Getenv("GOOGLE_API_KEYS"); multi != "" {
		for _, k := range strings.Split(multi, ",") {
			k = strings.TrimSpace(k)
			if k != "" {
				apiKeys = append(apiKeys, k)
			}
		}
	} else if single := os.Getenv("GOOGLE_API_KEY"); single != "" {
		apiKeys = []string{single}
	}

	return &Config{
		GoogleAPIKeys: apiKeys,
		GeminiAPIURL:  os.Getenv("GEMINI_API_URL"),
		PineconeKey:   os.Getenv("PINECONE_API_KEY"),
		PineconeHost:  os.Getenv("PINECONE_HOST"),
		Port:          port,
	}
}
