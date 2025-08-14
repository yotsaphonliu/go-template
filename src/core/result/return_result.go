package result

import (
	"net/http"
)

type Result struct {
	Code       int         `json:"code"`
	Message    string      `json:"message"`
	Pagination interface{} `json:"pagination,omitempty"`
	Data       interface{} `json:"data,omitempty"`
}

// Error error description
func (rs Result) Error() string {
	return rs.Message
}

// ErrorCode error code
func (rs Result) ErrorCode() int {
	return rs.Code
}

// HTTPStatusCode http status code
func (rs Result) HTTPStatusCode() int {
	switch rs.Code {
	case 0, 200: // success
		return http.StatusOK
	case 400: // bad request
		return http.StatusBadRequest
	case 404: // connection_error
		return http.StatusNotFound
	case 401: // unauthorized
		return http.StatusUnauthorized
	}

	return http.StatusInternalServerError
}

// NewResultWithMessage new result with message
func NewResultWithMessage(message string) Result {
	return Result{
		Message: message,
	}
}
