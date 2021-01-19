package index

// Document represents single document indexed
type Document interface {
	SetID(id int)
	GetID() int
	GetWords() []string
	GetData() string
	GetShortHTML(query string) string
	GetHighlightedHTML(query string) string
}

// Query represents single query for search
type Query interface {
	GetWords() []string
}

// SearchResult represents result of Search()
type SearchResult struct {
	Relevants  []Document
	TotalPages int
}

// GetResults represents result of Get()
type GetResults struct {
	Doc       Document
	NextID    int // set to -1 if not exists
	PrevID    int // set to -1 of not exists
	TotalDocs int
}
