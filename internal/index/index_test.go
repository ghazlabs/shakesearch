package index_test

import (
	"errors"
	"strings"
	"testing"

	"pulley.com/shakesearch/internal/errs"
	"pulley.com/shakesearch/internal/index"
)

// this standard stopwords from lucene index for english language
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
		&mockDocument{data: "hello world this is document 0"},
		&mockDocument{data: "cat is currently walking this is document 1 fish cat fly cat"},
		&mockDocument{data: "you know why we have this document this is document 2"},
		&mockDocument{data: "i'm fan of cat steven this is document 3"},
		&mockDocument{data: "perhaps i'm no longer exists maybe as cat guarding cat"},
	}
	// initialize index
	idx, err := index.New(
		index.Configs{
			Documents:     docs,
			ExcludedWords: stopWords,
			PageLimit:     2,
		},
	)
	if err != nil {
		t.Fatalf("unable to initialize index due: %v", err)
	}
	// search for `cat`, get second page (page 0)
	result, err := idx.Search(mockQuery("cat"), 0)
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
		gotID := result.Relevants[i].GetID()
		if gotID != expID {
			t.Fatalf("unexpected id, exp: %v, got: %v", expID, gotID)
		}
	}
	// search for `cat`, get second page (page 1)
	result, err = idx.Search(mockQuery("cat"), 1)
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
		gotID := result.Relevants[i].GetID()
		if gotID != expID {
			t.Fatalf("unexpected id, exp: %v, got: %v", expID, gotID)
		}
	}
}

// Test Page Not Found Search
func TestPageNotFoundSearch(t *testing.T) {
	// initialize documents to be indexed
	docs := []index.Document{
		&mockDocument{data: "hello world this is document 0"},
		&mockDocument{data: "cat is currently walking this is document 1 fish cat fly cat"},
		&mockDocument{data: "you know why we have this document this is document 2"},
		&mockDocument{data: "i'm fan of cat steven this is document 3"},
		&mockDocument{data: "perhaps i'm no longer exists maybe as cat guarding cat"},
	}
	// initialize index
	idx, err := index.New(
		index.Configs{
			Documents:     docs,
			ExcludedWords: stopWords,
			PageLimit:     2,
		},
	)
	if err != nil {
		t.Fatalf("unable to initialize index due: %v", err)
	}
	// search for `cat`, we expect to get only maximum 2 pages
	// but in here we fetch the forth page (page 3), it should
	// throw ERR_PAGE_NOT_FOUND
	_, err = idx.Search(mockQuery("cat"), 3)
	if errors.Is(err, errs.NewErrDocNotFound()) {
		t.Fatalf("unexpected error, exp: %v, got: %v", errs.NewErrDocNotFound(), err)
	}
}

// Test No Result Search
func TestNoResultSearch(t *testing.T) {
	// initialize documents to be indexed
	docs := []index.Document{
		&mockDocument{data: "hello world this is document 0"},
		&mockDocument{data: "cat is currently walking this is document 1 fish cat fly cat"},
		&mockDocument{data: "you know why we have this document this is document 2"},
		&mockDocument{data: "i'm fan of cat steven this is document 3"},
		&mockDocument{data: "perhaps i'm no longer exists maybe as cat guarding cat"},
	}
	// initialize index
	idx, err := index.New(
		index.Configs{
			Documents:     docs,
			ExcludedWords: stopWords,
			PageLimit:     2,
		},
	)
	if err != nil {
		t.Fatalf("unable to initialize index due: %v", err)
	}
	// search for `earth`, we should get empty result
	result, err := idx.Search(mockQuery("earth"), 0)
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

type mockQuery string

func (q mockQuery) GetWords() []string {
	return strings.Split(string(q), " ")
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
