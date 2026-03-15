package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/eyop23/insurance-go/config"
	"github.com/eyop23/insurance-go/models"
	"github.com/eyop23/insurance-go/services"
)

func main() {
	cfg := config.Load()

	data, err := os.ReadFile("data/ethiopian_music.json")
	if err != nil {
		log.Fatalf("Failed to read ethiopian_music.json: %v", err)
	}

	var artists []models.MusicArtist
	if err := json.Unmarshal(data, &artists); err != nil {
		log.Fatalf("Failed to parse JSON: %v", err)
	}

	fmt.Printf("Loaded %d Ethiopian music artists\n", len(artists))

	var vectors []models.PineconeVector
	batchSize := 10

	for i, artist := range artists {
		text := services.FlattenArtist(artist)
		embedding, err := services.GetEmbedding(text, cfg.GoogleAPIKey)
		if err != nil {
			fmt.Printf("Failed to embed %s: %v\n", artist.Name, err)
			continue
		}

		fmt.Printf("[%d/%d] Embedded %s (%d dims)\n", i+1, len(artists), artist.Name, len(embedding))

		vectors = append(vectors, models.PineconeVector{
			ID:     artist.ID,
			Values: embedding,
			Metadata: map[string]string{
				"text":        text,
				"name":        artist.Name,
				"genre":       strings.Join(artist.Genre, ", "),
				"instruments": strings.Join(artist.Instruments, ", "),
				"era":         artist.Era,
			},
		})

		if len(vectors) >= batchSize {
			if err := services.UpsertToPinecone(vectors, cfg.PineconeHost, cfg.PineconeKey); err != nil {
				log.Fatalf("Upsert failed: %v", err)
			}
			fmt.Printf("Upserted batch of %d vectors\n", len(vectors))
			vectors = nil
		}
	}

	if len(vectors) > 0 {
		if err := services.UpsertToPinecone(vectors, cfg.PineconeHost, cfg.PineconeKey); err != nil {
			log.Fatalf("Upsert failed: %v", err)
		}
		fmt.Printf("Upserted final batch of %d vectors\n", len(vectors))
	}

	fmt.Println("Done! All Ethiopian music artists seeded to Pinecone.")
}
