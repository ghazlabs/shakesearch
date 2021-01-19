package errs

import "net/http"

// Error represents error that will be used in this system
type Error struct {
	StatusCode int
	ErrCode    string
	Message    string
}

func (e *Error) Error() string {
	return e.ErrCode
}

// Is returns true if the error code is same
func (e *Error) Is(target error) bool {
	t, ok := target.(*Error)
	if !ok {
		return false
	}
	return e.ErrCode == t.ErrCode
}

// NewErrEmptyQuery returns new instance of ERR_EMPTY_QUERY
func NewErrEmptyQuery() *Error {
	return &Error{
		StatusCode: http.StatusBadRequest,
		ErrCode:    "ERR_EMPTY_QUERY",
	}
}

// NewErrPageNotFound returns new instance of ERR_PAGE_NOT_FOUND
func NewErrPageNotFound() *Error {
	return &Error{
		StatusCode: http.StatusNotFound,
		ErrCode:    "ERR_PAGE_NOT_FOUND",
	}
}

// NewErrDocNotFound returns new instance of ERR_DOC_NOT_FOUND
func NewErrDocNotFound() *Error {
	return &Error{
		StatusCode: http.StatusNotFound,
		ErrCode:    "ERR_DOC_NOT_FOUND",
	}
}

// NewErrInternalError returns new instance of ERR_INTERNAL_ERROR
func NewErrInternalError(err error) *Error {
	return &Error{
		StatusCode: http.StatusInternalServerError,
		ErrCode:    "ERR_INTERNAL_ERROR",
		Message:    err.Error(),
	}
}
