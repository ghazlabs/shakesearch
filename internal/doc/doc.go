package doc

import (
	"regexp"
	"sort"
	"strings"

	"pulley.com/shakesearch/internal/index"
	"pulley.com/shakesearch/internal/query"
)

// Document represents document that will be indexed
type Document struct {
	id            int
	data          string
	shortTag      Tag
	highlightTag  Tag
	lines         []string
	maxShortChars int
}

// Configs represents configs for Document
type Configs struct {
	Lines         []string
	ShortTag      Tag
	HighlightTag  Tag
	MaxShortChars int
}

// New returns new instance of Document
func New(c Configs) (*Document, error) {
	// construct documents
	doc := &Document{
		data:          strings.Join(c.Lines, "\n"),
		shortTag:      c.ShortTag,
		highlightTag:  c.HighlightTag,
		lines:         c.Lines,
		maxShortChars: c.MaxShortChars,
	}
	return doc, nil
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
//
// To construct this currently we do it like this:
//
// 1. break the query string into words
// 2. find the corresponding lines where these words appear
// 3. short the lines according to its index
// 4. merge the lines into one paragraph with following rules:
//		4.a. if current line is the next of previous line, append
//		     it immediately
//		4.b. if current line is not next of previous line, append
//			 "..." then append the current line
//		4.c. limit the resulted paragraph to 170 chars
// 5. wrap every words in paragraph with specified html tag
func (d *Document) GetShortHTML(words []string) string {
	// find the lines for each words
	lineMap := map[int]struct{}{}
	for _, word := range words {
		regex := buildMatchAnyCasePattern(word)
		for i := 0; i < len(d.lines); i++ {
			if !regex.MatchString(d.lines[i]) {
				continue
			}
			lineMap[i] = struct{}{}
		}
	}
	// if no line found, returns empty string
	if len(lineMap) == 0 {
		return ""
	}
	// convert line map into list
	lineIdxs := make([]int, 0, len(lineMap))
	for lineIdx := range lineMap {
		lineIdxs = append(lineIdxs, lineIdx)
	}
	// sort the list
	sort.Slice(lineIdxs, func(i, j int) bool {
		return lineIdxs[i] < lineIdxs[j]
	})
	// merge lines into paragraph
	pBuilder := &strings.Builder{}
	prevIdx := 0
	for i := 0; i < len(lineIdxs); i++ {
		currentIdx := lineIdxs[i]
		currentLine := strings.TrimSpace(d.lines[currentIdx])
		if i == 0 || currentIdx-1 == prevIdx {
			pBuilder.WriteString(currentLine)
		} else {
			pBuilder.WriteString("... " + currentLine)
		}
		if i < len(lineIdxs)-1 {
			pBuilder.WriteString(" ")
		}
		prevIdx = currentIdx
	}
	dataHTML := pBuilder.String()
	if len([]rune(dataHTML)) > d.maxShortChars {
		dataHTML = string([]rune(dataHTML)[:d.maxShortChars]) + "..."
	}
	// warp every words appears in paragraph by tag
	return warpWordsByTag(dataHTML, d.shortTag, words)
}

// GetHighlightedHTML returns full document data, but for
// every words in the query it will wrapped by specified
// highlight html tag
func (d *Document) GetHighlightedHTML(words []string) string {
	// make sure words are unique
	uniqueWords := query.Query(strings.Join(words, " ")).GetUniqueWords()
	// wrap every words with highlight tag
	return warpWordsByTag(d.data, d.highlightTag, uniqueWords)
}

func buildMatchAnyCasePattern(word string) *regexp.Regexp {
	patternBuilder := &strings.Builder{}
	for _, c := range word {
		patternBuilder.WriteString("(" + strings.ToLower(string(c)) + "|" + strings.ToUpper(string(c)) + ")")
	}
	return regexp.MustCompile(patternBuilder.String())
}

func warpWordsByTag(text string, tag Tag, words []string) string {
	dataHTML := text
	for _, word := range words {
		regex := buildMatchAnyCasePattern(word)
		foundWords := regex.FindAllString(dataHTML, -1)
		foundWordMap := map[string]struct{}{}
		for _, uw := range foundWords {
			foundWordMap[uw] = struct{}{}
		}
		for foundWord := range foundWordMap {
			dataHTML = strings.ReplaceAll(dataHTML, foundWord, tag.Start+foundWord+tag.End)
		}
	}
	return dataHTML
}

// GetLines returns document lines
func (d *Document) GetLines() []index.Line {
	lines := make([]index.Line, 0, len(d.lines))
	for _, lineStr := range d.lines {
		lines = append(lines, query.Query(lineStr))
	}
	return lines
}
