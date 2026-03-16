package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/eyop23/insurance-go/config"
	"github.com/eyop23/insurance-go/models"
	"github.com/eyop23/insurance-go/services"
)

func main() {
	resume := flag.Bool("resume", false, "Resume from where you left off (skips already embedded books)")
	flag.Parse()

	cfg := config.Load()

	data, err := os.ReadFile("data/books.json")
	if err != nil {
		log.Fatalf("Failed to read books.json: %v", err)
	}

	var books []models.Book
	if err := json.Unmarshal(data, &books); err != nil {
		log.Fatalf("Failed to parse JSON: %v", err)
	}

	fmt.Printf("Loaded %d books\n", len(books))

	startIndex := 0
	if *resume {
		count, err := services.GetPineconeVectorCount(cfg.PineconeHost, cfg.PineconeKey)
		if err != nil {
			log.Fatalf("Failed to get Pinecone vector count: %v", err)
		}
		startIndex = count
		fmt.Printf("Pinecone has %d vectors. Resuming from book #%d\n", count, startIndex+1)
	}

	if startIndex >= len(books) {
		fmt.Println("All books already seeded!")
		return
	}

	var vectors []models.PineconeVector
	batchSize := 10

	for i := startIndex; i < len(books); i++ {
		book := books[i]
		text := services.FlattenBook(book)
		embedding, err := services.GetEmbedding(text, cfg.GoogleAPIKey)
		if err != nil {
			fmt.Printf("[%d/%d] Failed to embed %s: %v\n", i+1, len(books), book.Title, err)
			fmt.Printf("Stopped at book #%d. Run again with --resume and a new API key to continue.\n", i+1)
			break
		}

		fmt.Printf("[%d/%d] Embedded %s (%d dims)\n", i+1, len(books), book.Title, len(embedding))

		vectors = append(vectors, models.PineconeVector{
			ID:     book.ID,
			Values: embedding,
			Metadata: map[string]string{
				"text":   text,
				"title":  book.Title,
				"genres": strings.Join(book.Genres, ", "),
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

	fmt.Println("Done! All books seeded to Pinecone.")
}
