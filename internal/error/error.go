package error

import "fmt"

type RapidLinksError struct {
	StatusCode int
	Message    string
}

func (e RapidLinksError) Error() string {
	return fmt.Sprintf("%d - %s", e.StatusCode, e.Message)
}

func (e RapidLinksError) GetStatusCode() int {
	return e.StatusCode
}

func NewRapidLinksError(message string, statusCode int) error {
	return RapidLinksError{
		StatusCode: statusCode,
		Message:    message,
	}
}
