package models

// --- Ethiopian Music ---

type Album struct {
	Title string `json:"title"`
	Year  int    `json:"year"`
	Label string `json:"label"`
}

type MusicArtist struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	BirthYear   int      `json:"birthYear"`
	DeathYear   int      `json:"deathYear,omitempty"`
	Origin      string   `json:"origin"`
	Genre       []string `json:"genre"`
	Instruments []string `json:"instruments"`
	Era         string   `json:"era"`
	Bio         string   `json:"bio"`
	Albums      []Album  `json:"albums"`
	FamousSongs []string `json:"famousSongs"`
	Awards      []string `json:"awards"`
	Influence   string   `json:"influence"`
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
