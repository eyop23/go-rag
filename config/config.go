package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	GoogleAPIKey string
	GeminiAPIURL string
	PineconeKey  string
	PineconeHost string
	Port         string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env not found, using system env")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8090"
	}

	return &Config{
		GoogleAPIKey: os.Getenv("GOOGLE_API_KEY"),
		GeminiAPIURL: os.Getenv("GEMINI_API_URL"),
		PineconeKey:  os.Getenv("PINECONE_API_KEY"),
		PineconeHost: os.Getenv("PINECONE_HOST"),
		Port:         port,
	}
}
