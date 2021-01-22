package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"time"

	"pulley.com/shakesearch/internal/doc"
	"pulley.com/shakesearch/internal/index"
	"pulley.com/shakesearch/internal/rest"
)

const (
	maxLinesInPage        = 50
	searchResultPageLimit = 10
)

const defaultPort = "3001"

var (
	shortTag = doc.Tag{
		Start: "<b>",
		End:   "</b>",
	}
	highlightTag = doc.Tag{
		Start: `<span style="highlight">"`,
		End:   `</span>`,
	}
	// this standard stopwords from lucene index for english language
	stopWords = []string{
		"a", "an", "and", "are", "as", "at", "be", "but", "by",
		"for", "if", "in", "into", "is", "it",
		"no", "not", "of", "on", "or", "such",
		"that", "the", "their", "then", "there", "these",
		"they", "this", "to", "was", "will", "with",
	}
)

func main() {
	// open txt file
	file, err := os.Open("completeworks.txt")
	if err != nil {
		log.Fatalf("unable to open source file due: %v", err)
	}
	// read the file line by line using scanner
	scanner := bufio.NewScanner(file)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	// break the lines into documents (a.k.a pages)
	countDocs := int(math.Ceil(float64(len(lines)) / float64(maxLinesInPage)))
	docs := make([]index.Document, 0, countDocs)
	for i := 0; i < countDocs; i++ {
		startIdx := i * maxLinesInPage
		endIdx := startIdx + maxLinesInPage
		if endIdx > countDocs {
			endIdx = countDocs
		}
		d, err := doc.New(doc.Configs{
			Lines:        lines[startIdx:endIdx],
			ShortTag:     shortTag,
			HighlightTag: highlightTag,
		})
		if err != nil {
			log.Fatalf("unable to create new doc due: %v", err)
		}
		docs = append(docs, d)
	}
	// initialize search index
	searchIdx, err := index.New(
		index.Configs{
			Documents:       docs,
			ExcludedWords:   stopWords,
			ResultPageLimit: searchResultPageLimit,
		},
	)
	if err != nil {
		log.Fatalf("unable to initialize search index due: %v", err)
	}
	// initialize api handler
	api := rest.New(searchIdx)
	// initialize http server
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	server := &http.Server{
		Addr:        fmt.Sprintf(":%v", defaultPort),
		Handler:     api.GetHandler(),
		ReadTimeout: 3 * time.Second,
	}
	// start http server
	log.Printf("server is listening on :%v", port)
	err = server.ListenAndServe()
	if err != nil {
		log.Fatalf("unable to start server due: %v", err)
	}
}
