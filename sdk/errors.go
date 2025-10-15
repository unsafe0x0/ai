package sdk

import "fmt"

type APIError struct {
	StatusCode int
	Message    string
	Body       []byte
}

func (e *APIError) Error() string {
	return fmt.Sprintf("APIError: %d - %s", e.StatusCode, e.Message)
}
