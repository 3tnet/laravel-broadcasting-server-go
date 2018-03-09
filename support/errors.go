package support

import (
	"encoding/json"
	"net/http"
)

// HTTPError represents an HTTP error with HTTP status code and error message
type HTTPError interface {
	error
	// StatusCode returns the HTTP status code of the error
	StatusCode() int
}

type httpError struct {
	statusCode int
	msg        error
}

func (e httpError) StatusCode() int {
	return e.statusCode
}

func (e httpError) Error() string {
	return e.msg.Error()
}

func NewHTTPError(statusCode int, msg error) HTTPError {
	return httpError{statusCode: statusCode, msg: msg}
}

// apiError represents an error that can be sent in an error response.
type APIError struct {
	// Status represents the HTTP status code
	status int `json:"-"`
	// ErrorCode is the code uniquely identifying an error
	// ErrorCode string `json:"error_code"`
	// Message is the error message that may be displayed to end users
	Message string `json:"message"`
	// Details specifies the additional error information
	Errors interface{} `json:"errors,omitempty"`
}

// Error returns the error message.
func (e APIError) Error() string {
	if jsonData, err := e.ToJson(); err == nil {
		return string(jsonData)
	} else {
		return err.Error()
	}
}

// StatusCode returns the HTTP status code.
func (e APIError) StatusCode() int {
	return e.status
}

func (e APIError) ToJson() ([]byte, error) {
	return json.Marshal(e)
}

func NewAPIError(status int, message string, errors ...interface{}) *APIError {
	apiError := &APIError{
		status:  status,
		Message: message,
		// ErrorCode: errorCode,
	}
	if len(errors) > 0 {
		apiError.Errors = errors[0]
	}
	return apiError
}

// InternalServerError creates a new API error representing an internal server error (HTTP 500)
func InternalServerError(message ...string) *APIError {
	var msg string
	if len(message) > 0 {
		msg = message[0]
	} else {
		msg = http.StatusText(http.StatusInternalServerError)
	}
	return NewAPIError(http.StatusInternalServerError, msg)
}

// NotFound creates a new API error representing a resource-not-found error (HTTP 404)
func NotFound(message ...string) *APIError {
	var msg string
	if len(message) > 0 {
		msg = message[0]
	} else {
		msg = http.StatusText(http.StatusNotFound)
	}
	return NewAPIError(http.StatusNotFound, msg)
}

// Unauthorized creates a new API error representing an authentication failure (HTTP 401)
func Unauthorized(message ...string) *APIError {
	var msg string
	if len(message) > 0 {
		msg = message[0]
	} else {
		msg = http.StatusText(http.StatusUnauthorized)
	}
	return NewAPIError(http.StatusUnauthorized, msg)
}

func Forbidden(message ...string) *APIError {
	var msg string
	if len(message) > 0 {
		msg = message[0]
	} else {
		msg = http.StatusText(http.StatusForbidden)
	}
	return NewAPIError(http.StatusForbidden, msg)
}
