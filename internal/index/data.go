package index

// Document represents single document indexed
type Document interface {
	SetID(id int)
	GetID() int
	GetWords() []string
	GetData() string
	GetShortHTML(query string) string
	GetHighlightedHTML(query string) string
	GetLines() []Line
}

// Query represents single query for search
type Query interface {
	GetWords() []string
	GetUniqueWords() []string
}

// Line represents single line in document
type Line Query

// SearchResult represents result of Search()
type SearchResult struct {
	Relevants  []Relevant
	TotalPages int
}

// Relevant represents single relevant result
type Relevant struct {
	Document   Document
	FoundWords []string
	Score      float64
}

// GetResults represents result of Get()
type GetResults struct {
	Doc       Document
	NextID    int // set to -1 if not exists
	PrevID    int // set to -1 of not exists
	TotalDocs int
}
