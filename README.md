# Book Explorer RAG — Go Backend

A Retrieval-Augmented Generation (RAG) API built with Go, Pinecone, and Google Gemini for querying a dataset of 2,032 books from Goodreads.

## Tech Stack

- **Go** + **Gin** — HTTP server
- **Pinecone** — Vector database (cosine similarity, 3072 dims)
- **Google Gemini** — Embedding (`gemini-embedding-001`) + LLM (`gemini-2.5-flash`)

## Project Structure

```
backend/
├── cmd/
│   ├── server/main.go      # API server entry point
│   └── seed/main.go         # Seed data into Pinecone (supports --resume)
├── config/config.go          # Environment config loader
├── models/types.go           # All data structs
├── services/
│   ├── embedding.go          # Gemini embedding generation
│   ├── pinecone.go           # Pinecone query, upsert & stats
│   ├── llm.go                # Gemini LLM answer generation
│   └── book.go               # Book record flattener
├── handlers/ask.go           # /ask endpoint handler
├── middleware/cors.go        # CORS middleware
├── data/
│   └── books.json            # 2,032 books with descriptions & genres
└── .env                      # API keys (not committed)
```

## Setup

1. Create a `.env` file:
```
GOOGLE_API_KEY=your-google-api-key
GEMINI_API_URL=https://generativelanguage.googleapis.com/v1/models/gemini-2.5-flash:generateContent
PINECONE_API_KEY=your-pinecone-api-key
PINECONE_HOST=https://your-index.svc.pinecone.io
PORT=8090
```

2. Install dependencies:
```bash
go mod tidy
```

3. Seed Pinecone with book data:
```bash
go run ./cmd/seed
```

The free tier of Gemini allows ~1,000 embeddings per API key per day. To seed all 2,032 books:
```bash
# First run — embeds ~1,000 books, then hits rate limit
go run ./cmd/seed

# Swap API key in .env, then resume from where it stopped
go run ./cmd/seed --resume
```

4. Start the server:
```bash
go run ./cmd/server
```

## API

### POST /ask

```json
{
  "query": "Recommend me a classic novel about justice"
}
```

Response:
```json
{
  "answer": "To Kill a Mockingbird by Harper Lee is a classic novel about justice..."
}
```

## How It Works

1. User query is converted to a 3072-dim vector via Gemini embedding
2. Pinecone finds the most similar book records via cosine similarity
3. Matched records are passed as context to Gemini LLM
4. LLM generates a grounded answer based only on the retrieved context
