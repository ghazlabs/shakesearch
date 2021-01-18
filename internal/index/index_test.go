package index_test

import (
	"strings"
	"testing"

	"pulley.com/shakesearch/internal/index"
)

var stopWords = []string{
	"a", "an", "and", "are", "as", "at", "be", "but", "by",
	"for", "if", "in", "into", "is", "it",
	"no", "not", "of", "on", "or", "such",
	"that", "the", "their", "then", "there", "these",
	"they", "this", "to", "was", "will", "with",
}

// Test Normal Search
func TestNormalSearch(t *testing.T) {
	// initialize documents to be indexed
	docs := []index.Document{
		&mockDocument{id: 0, data: "hello world this is document 0"},
		&mockDocument{id: 1, data: "cat is currently walking this is document 1 fish cat fly cat"},
		&mockDocument{id: 2, data: "you know why we have this document this is document 2"},
		&mockDocument{id: 3, data: "i'm fan of cat steven this is document 3"},
		&mockDocument{id: 4, data: "perhaps i'm no longer exists maybe as cat guarding cat"},
	}
	// initialize index
	idx, err := index.New(
		index.Configs{
			Documents:     docs,
			ExcludedWords: stopWords,
			PageLimit:     3,
		},
	)
	if err != nil {
		t.Fatalf("unable to initialize index due: %v", err)
	}
	// search for `cat`
	result, err := idx.Search(mockQuery("cat"), 0)
	if err != nil {
		t.Fatalf("unable to search due: %v", err)
	}
	// check result
	expIDs := []int{1, 4, 3}
	if len(result.Relevants) != len(expIDs) {
		t.Fatalf("unexpected relevant ids: %+v, exp: %v", result.Relevants, expIDs)
	}
	for i := 0; i < len(expIDs); i++ {
		expID := expIDs[i]
		gotID := result.Relevants[i].GetID()
		if gotID != expID {
			t.Fatalf("unexpected id, exp: %v, got: %v", expID, gotID)
		}
	}
}

// Test Page Not Found Search
func TestPageNotFoundSearch(t *testing.T) {
	// TODO
}

// Test No Result Search
func TestNoResultSearch(t *testing.T) {
	// TODO
}

type mockQuery string

func (q mockQuery) GetWords() []string {
	return strings.Split(string(q), " ")
}

type mockDocument struct {
	id   int
	data string
}

func (d *mockDocument) GetID() int {
	return d.id
}

func (d *mockDocument) GetWords() []string {
	return strings.Split(d.data, " ")
}

func (d *mockDocument) GetData() string {
	return d.data
}

func (d *mockDocument) GetShortHTML(query string) string {
	// we don't use this in search, so it's fine to just
	// returns arbitrary string
	return d.data
}

func (d *mockDocument) GetHighlightedHTML(query string) string {
	// we don't use this in search, so it's fine to just
	// returns arbitrary string
	return d.data
}
