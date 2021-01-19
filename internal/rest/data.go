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
	Data    interface{} `json:"data.omitempty"`
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
	TotalPages  int              `json:"total_pages"`
}

func newSearchData(queryString string, currentPage int, result *index.SearchResult) *searchData {
	searchRelevants := make([]searchRelevant, 0, len(result.Relevants))
	for _, doc := range result.Relevants {
		searchRelevants = append(searchRelevants, searchRelevant{
			ShortHTML: doc.GetShortHTML(queryString),
			URL:       fmt.Sprintf("/pages/%v?q=%v", doc.GetID(), queryString),
		})
	}
	return &searchData{
		Relevants:   searchRelevants,
		CurrentPage: currentPage,
		TotalPages:  result.TotalPages,
	}
}

type searchRelevant struct {
	ShortHTML string `json:"short_html"`
	URL       string `json:"url"`
}

type viewPageData struct {
	BodyHTML    string `json:"body_html"`
	CurrentPage int    `json:"current_page"`
	PrevPage    int    `json:"prev_page"`
	NextPage    int    `json:"next_page"`
	TotalPages  int    `json:"total_pages"`
}
