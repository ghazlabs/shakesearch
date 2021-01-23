package query

import (
	"regexp"
	"strings"
)

// regex for finding words, in here we use non whitespace characters
var rgx = regexp.MustCompile(`\S+`)

// replacer for weird symbols
var rplcr = strings.NewReplacer("’", "'", "’", "'")

// Query represents query for search
type Query string

// GetWords returns every words in query, case sensitive
func (q Query) GetWords() []string {
	words := rgx.FindAllString(string(q), -1)
	filteredWords := make([]string, 0, len(words))
	for _, word := range words {
		word = rplcr.Replace((strings.Trim(word, ".,*:][#;!?")))
		if len(word) == 0 {
			continue
		}
		filteredWords = append(filteredWords, word)
	}
	return filteredWords
}

// GetUniqueWords returns list of unique words
func (q Query) GetUniqueWords() []string {
	wordMap := map[string]struct{}{}
	words := q.GetWords()
	for _, word := range words {
		wordMap[word] = struct{}{}
	}
	uniqueWords := make([]string, 0, len(words))
	for word := range wordMap {
		uniqueWords = append(uniqueWords, word)
	}
	return uniqueWords
}
