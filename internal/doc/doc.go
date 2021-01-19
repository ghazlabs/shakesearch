package doc

import (
	"strings"

	"pulley.com/shakesearch/internal/query"
)

// Document represents document that will be indexed
type Document struct {
	id           int
	data         string
	shortTag     Tag
	highlightTag Tag
	wordMap      map[string][]WordPos
}

// Configs represents configs for Document
type Configs struct {
	Data         string
	ShortTag     Tag
	HighlightTag Tag
}

// New returns new instance of Document
func New(c Configs) (*Document, error) {
	// TODO
	return nil, nil
}

// SetID is used for setting the document id
func (d *Document) SetID(id int) {
	d.id = id
}

// GetID returns document id
func (d *Document) GetID() int {
	return d.id
}

// GetWords returns words in document data
func (d *Document) GetWords() []string {
	return query.Query(d.data).GetWords()
}

// GetData returns document data
func (d *Document) GetData() string {
	return d.data
}

// GetShortHTML returns short document data, for purpose of
// quick view based on given query. Think it like search result
// body when we search something on google.
func (d *Document) GetShortHTML(queryString string) string {
	// TODO
	return ""
}

// GetHighlightedHTML returns full document data, but for
// every words in the query it will wrapped by specified
// highlight html tag
func (d *Document) GetHighlightedHTML(queryString string) string {
	// breakdown query into words
	words := query.Query(queryString).GetWords()
	// wrap every words with highlight tag
	dataHTML := d.data
	for _, word := range words {
		dataHTML = strings.ReplaceAll(dataHTML, word, d.highlightTag.Start+word+d.highlightTag.End)
	}
	return dataHTML
}
