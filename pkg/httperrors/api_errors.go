package httperrors

import (
	"fmt"
	"net/http"
)

type APIErrorType string

const (
	ErrNotFound     APIErrorType = "NOT_FOUND"
	ErrConflict     APIErrorType = "CONFLICT"
	ErrBadRequest   APIErrorType = "BAD_REQUEST"
	ErrInternal     APIErrorType = "INTERNAL_ERROR"
	ErrValidation   APIErrorType = "VALIDATION_ERROR"
	ErrUnauthorized APIErrorType = "UNAUTHORIZED"
	ErrTimeout      APIErrorType = "TIMEOUT"
	ErrUnavailable  APIErrorType = "SERVICE_UNAVAILABLE"
	ErrForbidden    APIErrorType = "FORBIDDEN"
)

// APIError represents the structure of an HTTP error for APIs.
type APIError struct {
	Type    APIErrorType   `json:"type"`
	Code    int            `json:"code"`
	Message string         `json:"message"`
	Details string         `json:"details,omitempty"`
	Context map[string]any `json:"context,omitempty"`
}

// Error implements the error interface.
func (e *APIError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("%s: %s (%s)", e.Type, e.Message, e.Details)
	}
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

// HTTP status code mapping for each error type.
var httpStatus = map[APIErrorType]int{
	ErrBadRequest:   http.StatusBadRequest,
	ErrNotFound:     http.StatusNotFound,
	ErrConflict:     http.StatusConflict,
	ErrInternal:     http.StatusInternalServerError,
	ErrValidation:   http.StatusBadRequest,
	ErrUnauthorized: http.StatusUnauthorized,
	ErrTimeout:      http.StatusGatewayTimeout,
	ErrUnavailable:  http.StatusServiceUnavailable,
	ErrForbidden:    http.StatusForbidden,
}

// New creates a new APIError with the given type, message, and optional details/context.
func New(errType APIErrorType, message string, details string, context map[string]any) *APIError {
	code, ok := httpStatus[errType]
	if !ok {
		code = http.StatusInternalServerError
		errType = ErrInternal
	}

	return &APIError{
		Type:    errType,
		Code:    code,
		Message: message,
		Details: details,
		Context: context,
	}
}

// NewSimple creates a new APIError with only type and message.
func NewSimple(errType APIErrorType, message string) *APIError {
	return New(errType, message, "", nil)
}

// StatusCode returns the HTTP status code for the error.
func (e *APIError) StatusCode() int {
	return e.Code
}
