package errors

import "fmt"

type APIError struct {
	Message string `json:"message"`
}

func (e *APIError) Error() string {
	return e.Message
}

func NewAPIError(message string) *APIError {
	return &APIError{
		Message: message,
	}
}

func NewAPIErrorf(format string, a ...any) *APIError {
	return &APIError{
		Message: fmt.Sprintf(format, a...),
	}
}
