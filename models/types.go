package models

// --- Books ---

type Book struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Genres      []string `json:"genres"`
	URL         string   `json:"url"`
}

// --- API ---

type AskRequest struct {
	Query string `json:"query"`
}

type AskResponse struct {
	Answer string `json:"answer"`
	Error  string `json:"error,omitempty"`
}

// --- Gemini ---

type GeminiEmbedRequest struct {
	Content struct {
		Parts []struct {
			Text string `json:"text"`
		} `json:"parts"`
	} `json:"content"`
}

type GeminiEmbedResponse struct {
	Embedding struct {
		Values []float64 `json:"values"`
	} `json:"embedding"`
}

type GeminiGenRequest struct {
	Contents []struct {
		Parts []struct {
			Text string `json:"text"`
		} `json:"parts"`
	} `json:"contents"`
}

type GeminiGenResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

// --- Pinecone ---

type PineconeVector struct {
	ID       string            `json:"id"`
	Values   []float64         `json:"values"`
	Metadata map[string]string `json:"metadata"`
}

type PineconeUpsertRequest struct {
	Vectors []PineconeVector `json:"vectors"`
}

type PineconeQueryRequest struct {
	Vector          []float64 `json:"vector"`
	TopK            int       `json:"topK"`
	IncludeMetadata bool      `json:"includeMetadata"`
}

type PineconeMatch struct {
	ID       string            `json:"id"`
	Score    float64           `json:"score"`
	Metadata map[string]string `json:"metadata"`
}

type PineconeQueryResponse struct {
	Matches []PineconeMatch `json:"matches"`
}
