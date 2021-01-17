package index

import "errors"

// ErrDocNotFound represents error returned when document is not found
var ErrDocNotFound = errors.New("ERR_DOC_NOT_FOUND")

// Index represents in-memory index for documents
type Index struct {
	docMap         map[int]Document
	revIndexMap    map[string][]int
	excludeWordMap map[string]struct{}
}

// Configs holds configs for Index
type Configs struct {
	Documents     []Document
	ExcludedWords []string
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

// Search returns relevant documents based on given query
func (i *Index) Search(q string, page int) (*SearchResult, error) {
	// TODO
	return nil, nil
}

// Get returns index for given id, throws error if id is not exists
func (i *Index) Get(id int) (*Document, error) {
	v, ok := i.docMap[id]
	if !ok {
		return nil, ErrDocNotFound
	}
	return &v, nil
}
