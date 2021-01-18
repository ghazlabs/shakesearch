package index

import (
	"errors"
	"math"
	"sort"
)

var (
	// ErrDocNotFound represents error returned when document is not found
	ErrDocNotFound = errors.New("ERR_DOC_NOT_FOUND")
	// ErrPageNotFound represents error returned when page is not found
	ErrPageNotFound = errors.New("ERR_PAGE_NOT_FOUND")
)

// Index represents in-memory index for documents
type Index struct {
	docMap         map[int]Document
	revIndexMap    map[string][]int
	excludeWordMap map[string]struct{}
	pageLimit      int
}

// Configs holds configs for Index
type Configs struct {
	Documents     []Document
	ExcludedWords []string
	PageLimit     int
}

// New returns new instance of Index
func New(c Configs) (*Index, error) {
	// construct exclude word map
	excludeWordMap := map[string]struct{}{}
	for _, word := range c.ExcludedWords {
		excludeWordMap[word] = struct{}{}
	}
	// construct document map & reverse index map
	docMap := map[int]Document{}
	revIndexMap := map[string][]int{}
	for _, doc := range c.Documents {
		// set the document in document map
		docMap[doc.GetID()] = doc
		// get doc words and iterate on them
		for _, word := range doc.GetWords() {
			// if word is excluded, just skip it
			_, skipped := excludeWordMap[word]
			if skipped {
				continue
			}
			// get current document ids for word
			v, ok := revIndexMap[word]
			if !ok {
				// if it's new word, initialize new list
				v = []int{}
			}
			// insert current document id to map
			v = append(v, doc.GetID())
			revIndexMap[word] = v
		}
	}
	// initialize index
	i := Index{
		docMap:         docMap,
		revIndexMap:    revIndexMap,
		excludeWordMap: excludeWordMap,
	}
	return &i, nil
}

// Search returns relevant documents based on given query, page start from 0
func (i *Index) Search(q Query, page int) (*SearchResult, error) {
	// break query into words
	words := q.GetWords()
	// create appearance counter map
	counterMap := map[int]int{}
	// get document ids for each query words
	for _, word := range words {
		// if word is excluded word just continue
		_, excluded := i.excludeWordMap[word]
		if excluded {
			continue
		}
		// increment appearance counter map
		docIDs := i.revIndexMap[word]
		for _, docID := range docIDs {
			counterMap[docID]++
		}
	}
	// convert appearance counter map to list of tupple (docID, counter)
	tupples := make([]tuppleDocIDCounter, 0, len(counterMap))
	for docID, counter := range counterMap {
		tupples = append(tupples, tuppleDocIDCounter{
			DocID:   docID,
			Counter: counter,
		})
	}
	// sort the list from highest to lowest appearance
	sort.Slice(tupples, func(i, j int) bool {
		return tupples[i].Counter < tupples[j].Counter
	})
	// check if page is exists
	totalPages := int(math.Ceil(float64(len(tupples)) / float64(i.pageLimit)))
	if page >= totalPages {
		return nil, ErrPageNotFound
	}
	// select the elements for given page
	startIdx := page * i.pageLimit
	endIdx := startIdx + i.pageLimit
	if endIdx >= len(tupples) {
		endIdx = len(tupples)
	}
	// convert back tupples to list of docs
	selectedTupples := tupples[startIdx:endIdx]
	var docs []Document
	for _, tupple := range selectedTupples {
		docs = append(docs, i.docMap[tupple.DocID])
	}
	result := &SearchResult{
		Relevants:  docs,
		TotalPages: totalPages,
	}
	return result, nil
}

type tuppleDocIDCounter struct {
	DocID   int
	Counter int
}

// Get returns index for given id, throws error if id is not exists
func (i *Index) Get(id int) (*Document, error) {
	v, ok := i.docMap[id]
	if !ok {
		return nil, ErrDocNotFound
	}
	return &v, nil
}
