package rest

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"pulley.com/shakesearch/internal/errs"
	"pulley.com/shakesearch/internal/index"
	"pulley.com/shakesearch/internal/query"
)

// API is used for handling http api requests
type API struct {
	idx index.Index
}

// GetHandler returns http handler for API
func (a *API) GetHandler() http.Handler {
	r := chi.NewRouter()
	r.Get("/search", a.serveSearch)
	r.Get("/pages/{page_number}", a.serveViewPage)

	return r
}

const defPageNumber = 1

func (a *API) serveSearch(w http.ResponseWriter, r *http.Request) {
	// get query string from query params
	queryString := r.URL.Query().Get("q")
	if len(queryString) == 0 {
		render.Render(w, r, newErrResp(errs.NewErrEmptyQuery()))
		return
	}
	// get page number from query params if any
	pageNumber, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if pageNumber == 0 {
		pageNumber = defPageNumber
	}
	// execute search in index, we decrease the page number by 1
	// because the page start with 0 in index
	result, err := a.idx.Search(query.Query(queryString), pageNumber-1)
	if err != nil {
		render.Render(w, r, newErrResp(err))
		return
	}
	// output success response
	render.Render(w, r, newSuccessResp(newSearchData(queryString, pageNumber, result)))
}

func (a *API) serveViewPage(w http.ResponseWriter, r *http.Request) {
	// TODO
}
