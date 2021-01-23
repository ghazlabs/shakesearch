package index

import (
	"math"
	"sort"
	"strings"

	"pulley.com/shakesearch/internal/errs"
)

// Index represents in-memory index for documents
type Index struct {
	docs            []Document
	revIndexWordMap map[string][]int
	excludeWordMap  map[string]struct{}
	resultPageLimit int
}

// Configs holds configs for Index
type Configs struct {
	Documents       []Document
	ExcludedWords   []string
	ResultPageLimit int
}

// New returns new instance of Index
func New(configs Configs) (*Index, error) {
	// construct exclude word map
	excludeWordMap := map[string]struct{}{}
	for _, word := range configs.ExcludedWords {
		excludeWordMap[word] = struct{}{}
	}
	// construct document map & reverse index map
	revIndexMap := map[string][]int{}
	for i := 0; i < len(configs.Documents); i++ {
		// set doc id to value of i
		configs.Documents[i].SetID(i)
		// get doc lines and iterate on them
		for _, line := range configs.Documents[i].GetLines() {
			// get unique words in the line
			words := line.GetWords()
			// skip line that has less than 3 words, because most likely it doesn't
			// contain any context
			if len(words) < 3 {
				continue
			}
			for _, word := range words {
				// if word is all capitals, most likely it's part of start dialog
				// marking, in which means most likely have no context
				if word == strings.ToUpper(word) {
					continue
				}
				// lower word case, this is to make the search process
				// case insensitive
				word = strings.ToLower(word)
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
				v = append(v, configs.Documents[i].GetID())
				revIndexMap[word] = v
			}
		}
	}
	// initialize index
	i := Index{
		docs:            configs.Documents,
		revIndexWordMap: revIndexMap,
		excludeWordMap:  excludeWordMap,
		resultPageLimit: configs.ResultPageLimit,
	}
	return &i, nil
}

// Search returns relevant documents based on given query, page start from 0
func (i *Index) Search(q Query, page int) (*SearchResult, error) {
	// break query into words
	words := q.GetUniqueWords()
	// create appearance counter map, this is used for document scoring
	counterMap := map[int]map[string]int{}
	// get document ids for each query words
	for _, word := range words {
		// lowercase the word, so we could do insensitive
		// case search
		word = strings.ToLower(word)
		// if word is excluded word just continue
		_, excluded := i.excludeWordMap[word]
		if excluded {
			continue
		}
		// increment appearance counter map
		docIDs := i.revIndexWordMap[word]
		for _, docID := range docIDs {
			v, ok := counterMap[docID]
			if !ok {
				v = map[string]int{}
			}
			v[word]++
			counterMap[docID] = v
		}
	}
	// convert appearance counter map to list of Relevant
	relevants := make([]Relevant, 0, len(counterMap))
	for docID, wordCounterMap := range counterMap {
		// doc_score = total_count * total_words^3
		totalCount := 0
		totalWords := len(wordCounterMap)
		foundWords := make([]string, 0, len(wordCounterMap))
		for word, count := range wordCounterMap {
			totalCount += count
			foundWords = append(foundWords, word)
		}
		score := float64(totalCount) * math.Pow(float64(totalWords), 3)
		relevants = append(relevants, Relevant{
			Document:   i.docs[docID],
			FoundWords: foundWords,
			Score:      score,
		})
	}
	// sort relevants list from highest to lowest score
	sort.Slice(relevants, func(i, j int) bool {
		// if score is same, sort by document id, make the earlier
		// document higher order
		prevScore := relevants[i].Score
		nextScore := relevants[j].Score
		if prevScore == nextScore {
			return relevants[i].Document.GetID() < relevants[j].Document.GetID()
		}
		// sort by score descendingly
		return prevScore > nextScore
	})
	// if totalPages is zero, returns empty list and set total pages to 1
	// notice that we set the total pages into 1 because logically we want
	// to show the empty result in first page
	totalPages := int(math.Ceil(float64(len(relevants)) / float64(i.resultPageLimit)))
	if totalPages == 0 && page == 0 {
		return &SearchResult{Relevants: nil, TotalPages: 1}, nil
	}
	// if page is not exist returns error
	if page >= totalPages {
		return nil, errs.NewErrPageNotFound()
	}
	// select the elements for given page
	startIdx := page * i.resultPageLimit
	endIdx := startIdx + i.resultPageLimit
	if endIdx >= len(relevants) {
		endIdx = len(relevants)
	}
	// convert back tupples to list of docs
	result := &SearchResult{
		Relevants:  relevants[startIdx:endIdx],
		TotalPages: totalPages,
	}
	return result, nil
}

// Get returns index for given id, throws error if id is not exists
func (i *Index) Get(id int) (*GetResults, error) {
	if id >= len(i.docs) {
		return nil, errs.NewErrDocNotFound()
	}
	// set next id
	nextID := -1
	if id < len(i.docs)-1 {
		nextID = id + 1
	}
	// set prev id
	prevID := -1
	if id > 0 {
		prevID = id - 1
	}
	// construct result
	result := &GetResults{
		Doc:       i.docs[id],
		NextID:    nextID,
		PrevID:    prevID,
		TotalDocs: len(i.docs),
	}
	return result, nil
}
