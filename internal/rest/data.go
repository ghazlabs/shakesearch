package rest

import (
	"fmt"
	"log"
	"net/http"

	"pulley.com/shakesearch/internal/errs"
	"pulley.com/shakesearch/internal/index"
)

type apiResp struct {
	Status  int         `json:"-"`
	OK      bool        `json:"ok"`
	Data    interface{} `json:"data,omitempty"`
	ErrCode string      `json:"err,omitempty"`
}

func (rs *apiResp) Render(w http.ResponseWriter, r *http.Request) error {
	// we do nothing here, because we just want to implement Renderer
	// interface in github.com/go-chi/render
	return nil
}

func newSuccessResp(data interface{}) *apiResp {
	return &apiResp{
		Status: http.StatusOK,
		OK:     true,
		Data:   data,
	}
}

func newErrResp(err error) *apiResp {
	e, ok := err.(*errs.Error)
	if !ok {
		log.Printf("[RDebug] internal error: %v", err)
		e = errs.NewErrInternalError(err)
	}
	return &apiResp{
		Status:  e.StatusCode,
		OK:      false,
		ErrCode: e.ErrCode,
	}
}

type searchData struct {
	Relevants   []searchRelevant `json:"relevants"`
	CurrentPage int              `json:"current_page"`
	PrevPage    *int             `json:"prev_page"`
	NextPage    *int             `json:"next_page"`
	TotalPages  int              `json:"total_pages"`
}

func newSearchData(queryString string, currentPage int, result *index.SearchResult) *searchData {
	searchRelevants := make([]searchRelevant, 0, len(result.Relevants))
	for _, doc := range result.Relevants {
		searchRelevants = append(searchRelevants, searchRelevant{
			ID:        doc.GetID() + 1,
			Title:     fmt.Sprintf("Page %v", doc.GetID()+1),
			ShortHTML: doc.GetShortHTML(queryString),
		})
	}
	d := &searchData{
		Relevants:   searchRelevants,
		CurrentPage: currentPage,
		TotalPages:  result.TotalPages,
	}
	if currentPage != 1 {
		prevPage := currentPage - 1
		d.PrevPage = &prevPage
	}
	if currentPage != result.TotalPages {
		nextPage := currentPage + 1
		d.NextPage = &nextPage
	}
	return d
}

type searchRelevant struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	ShortHTML string `json:"short_html"`
}

type viewPageData struct {
	BodyHTML    string `json:"body_html"`
	CurrentPage int    `json:"current_page"`
	PrevPage    *int   `json:"prev_page"`
	NextPage    *int   `json:"next_page"`
	TotalPages  int    `json:"total_pages"`
}

func newViewPageData(queryString string, currentPage int, result *index.GetResults) *viewPageData {
	d := &viewPageData{
		BodyHTML:    result.Doc.GetHighlightedHTML(queryString),
		CurrentPage: currentPage,
		TotalPages:  result.TotalDocs,
	}
	// set next page
	if result.NextID != -1 {
		nextPage := result.NextID + 1 // +1 because page started with 1, yet id with 0
		d.NextPage = &nextPage
	}
	// set prev page
	if result.PrevID != -1 {
		prevPage := result.PrevID + 1 // +1 because page started with 1, yet id with 0
		d.PrevPage = &prevPage
	}
	return d
}
