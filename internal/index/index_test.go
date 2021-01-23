package index_test

import (
	"errors"
	"strings"
	"testing"

	"pulley.com/shakesearch/internal/errs"
	"pulley.com/shakesearch/internal/index"
	"pulley.com/shakesearch/internal/query"
)

// this standard stopwords from lucene index for english language
var stopWords = []string{
	"a", "an", "and", "are", "as", "at", "be", "but", "by",
	"for", "if", "in", "into", "is", "it",
	"no", "not", "of", "on", "or", "such",
	"that", "the", "their", "then", "there", "these",
	"they", "this", "to", "was", "will", "with",
}

// Test Single Word Search
func TestSingleWordSearch(t *testing.T) {
	// initialize documents to be indexed
	docs := []index.Document{
		&mockDocument{data: "Hello, World! This is document 0."},
		&mockDocument{data: "Cat is currently walking! This is document 1! Fish cat fly cat."},
		&mockDocument{data: "You know why we have this document? This is document 2."},
		&mockDocument{data: "I'm fan of Cat Steven! This is document 3."},
		&mockDocument{data: "Perhaps I'm no longer exists maybe as cat guarding cat. This is document 4."},
	}
	// initialize index
	idx, err := index.New(
		index.Configs{
			Documents:       docs,
			ExcludedWords:   stopWords,
			ResultPageLimit: 2,
		},
	)
	if err != nil {
		t.Fatalf("unable to initialize index due: %v", err)
	}
	// search for `cat`, get second page (page 0)
	result, err := idx.Search(query.Query("cat"), 0)
	if err != nil {
		t.Fatalf("unable to search due: %v", err)
	}
	// check total pages, we should get 2 total pages from result
	expTotalPages := 2
	if result.TotalPages != expTotalPages {
		t.Fatalf("unexpected total pages, exp: %v, got: %v", expTotalPages, result.TotalPages)
	}
	// check result on first page
	expIDs := []int{1, 4}
	if len(result.Relevants) != len(expIDs) {
		t.Fatalf("unexpected relevant ids: %+v, exp: %v", result.Relevants, expIDs)
	}
	for i := 0; i < len(expIDs); i++ {
		expID := expIDs[i]
		gotID := result.Relevants[i].Document.GetID()
		if gotID != expID {
			t.Fatalf("unexpected id, exp: %v, got: %v", expID, gotID)
		}
	}
	// search for `cat`, get second page (page 1)
	result, err = idx.Search(query.Query("cat"), 1)
	if err != nil {
		t.Fatalf("unable to search due: %v", err)
	}
	// check total pages, we should get 2 total pages from result (no different when we get first page)
	if result.TotalPages != expTotalPages {
		t.Fatalf("unexpected total pages, exp: %v, got: %v", expTotalPages, result.TotalPages)
	}
	// check result on second page
	expIDs = []int{3}
	if len(result.Relevants) != len(expIDs) {
		t.Fatalf("unexpected relevant ids: %+v, exp: %v", result.Relevants, expIDs)
	}
	for i := 0; i < len(expIDs); i++ {
		expID := expIDs[i]
		gotID := result.Relevants[i].Document.GetID()
		if gotID != expID {
			t.Fatalf("unexpected id, exp: %v, got: %v", expID, gotID)
		}
	}
}

// Test Multiple Words Search
func TestMultipleWordsSearch(t *testing.T) {
	// initialize documents to be indexed
	docs := []index.Document{
		&mockDocument{data: "Hello, World! This is document 0."},
		&mockDocument{data: "Cat is currently walking! This is document 1! Fish cat fly cat."},
		&mockDocument{data: "You know why we have this document? This is document 2."},
		&mockDocument{data: "I'm fan of Cat Steven! This is document 3."},
		&mockDocument{data: "Perhaps I'm no longer exists maybe as cat guarding cat. This is document 4."},
	}
	// initialize index
	idx, err := index.New(
		index.Configs{
			Documents:       docs,
			ExcludedWords:   stopWords,
			ResultPageLimit: 2,
		},
	)
	if err != nil {
		t.Fatalf("unable to initialize index due: %v", err)
	}
	// search for `cat`, get second page (page 0)
	result, err := idx.Search(query.Query("Cat Steven"), 0)
	if err != nil {
		t.Fatalf("unable to search due: %v", err)
	}
	// check total pages, we should get 2 total pages from result
	expTotalPages := 2
	if result.TotalPages != expTotalPages {
		t.Fatalf("unexpected total pages, exp: %v, got: %v", expTotalPages, result.TotalPages)
	}
	// check result on first page
	expIDs := []int{3, 1}
	if len(result.Relevants) != len(expIDs) {
		t.Fatalf("unexpected relevant ids: %+v, exp: %v", result.Relevants, expIDs)
	}
	for i := 0; i < len(expIDs); i++ {
		expID := expIDs[i]
		gotID := result.Relevants[i].Document.GetID()
		if gotID != expID {
			t.Fatalf("unexpected id, exp: %v, got: %v", expID, gotID)
		}
	}
	// search for `cat`, get second page (page 1)
	result, err = idx.Search(query.Query("Cat Steven"), 1)
	if err != nil {
		t.Fatalf("unable to search due: %v", err)
	}
	// check total pages, we should get 2 total pages from result (no different when we get first page)
	if result.TotalPages != expTotalPages {
		t.Fatalf("unexpected total pages, exp: %v, got: %v", expTotalPages, result.TotalPages)
	}
	// check result on second page
	expIDs = []int{4}
	if len(result.Relevants) != len(expIDs) {
		t.Fatalf("unexpected relevant ids: %+v, exp: %v", result.Relevants, expIDs)
	}
	for i := 0; i < len(expIDs); i++ {
		expID := expIDs[i]
		gotID := result.Relevants[i].Document.GetID()
		if gotID != expID {
			t.Fatalf("unexpected id, exp: %v, got: %v", expID, gotID)
		}
	}
}

// Test Page Not Found Search
func TestPageNotFoundSearch(t *testing.T) {
	// initialize documents to be indexed
	docs := []index.Document{
		&mockDocument{data: "Hello, World! This is document 0."},
		&mockDocument{data: "Cat is currently walking! This is document 1! Fish cat fly cat."},
		&mockDocument{data: "You know why we have this document? This is document 2."},
		&mockDocument{data: "I'm fan of Cat Steven! This is document 3."},
		&mockDocument{data: "Perhaps I'm no longer exists maybe as cat guarding cat. This is document 4."},
	}
	// initialize index
	idx, err := index.New(
		index.Configs{
			Documents:       docs,
			ExcludedWords:   stopWords,
			ResultPageLimit: 2,
		},
	)
	if err != nil {
		t.Fatalf("unable to initialize index due: %v", err)
	}
	// search for `cat`, we expect to get only maximum 2 pages
	// but in here we fetch the forth page (page 3), it should
	// throw ERR_PAGE_NOT_FOUND
	_, err = idx.Search(query.Query("cat"), 3)
	if errors.Is(err, errs.NewErrDocNotFound()) {
		t.Fatalf("unexpected error, exp: %v, got: %v", errs.NewErrDocNotFound(), err)
	}
}

// Test No Result Search
func TestNoResultSearch(t *testing.T) {
	// initialize documents to be indexed
	docs := []index.Document{
		&mockDocument{data: "Hello, World! This is document 0."},
		&mockDocument{data: "Cat is currently walking! This is document 1! Fish cat fly cat."},
		&mockDocument{data: "You know why we have this document? This is document 2."},
		&mockDocument{data: "I'm fan of Cat Steven! This is document 3."},
		&mockDocument{data: "Perhaps I'm no longer exists maybe as cat guarding cat. This is document 4."},
	}
	// initialize index
	idx, err := index.New(
		index.Configs{
			Documents:       docs,
			ExcludedWords:   stopWords,
			ResultPageLimit: 2,
		},
	)
	if err != nil {
		t.Fatalf("unable to initialize index due: %v", err)
	}
	// search for `earth`, we should get empty result
	result, err := idx.Search(query.Query("earth"), 0)
	if err != nil {
		t.Fatalf("unable to search due: %v", err)
	}
	// check result
	expResultIDs := []int{}
	if len(result.Relevants) != len(expResultIDs) {
		t.Fatalf("search result should be empty, got: %v", result.Relevants)
	}
	// returned total pages should be 1
	expTotalPages := 1
	if result.TotalPages != expTotalPages {
		t.Fatalf("unexpected total pages, exp: %v, got: %v", expTotalPages, result.TotalPages)
	}
}

type mockDocument struct {
	id   int
	data string
}

func (d *mockDocument) SetID(id int) {
	d.id = id
}

func (d *mockDocument) GetID() int {
	return d.id
}

func (d *mockDocument) GetWords() []string {
	return query.Query(d.data).GetWords()
}

func (d *mockDocument) GetData() string {
	return d.data
}

func (d *mockDocument) GetShortHTML(words []string) string {
	// we don't use this in search, so it's fine to just
	// returns arbitrary string
	return d.data
}

func (d *mockDocument) GetHighlightedHTML(words []string) string {
	// we don't use this in search, so it's fine to just
	// returns arbitrary string
	return d.data
}

func (d *mockDocument) GetLines() []index.Line {
	lineStrs := strings.Split(d.data, "\n")
	lines := make([]index.Line, 0, len(lineStrs))
	for _, lineStr := range lineStrs {
		lines = append(lines, query.Query(lineStr))
	}
	return lines
}
