package index

// Document represents single document indexed
type Document interface {
	GetID() int
	GetWords() []string
	GetData() string
	GetShortHTML() string
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
