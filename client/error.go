package client

import "fmt"

//Error represents a Generic API Client error
type Error struct {
	error
	URL            string
	HTTPStatusCode int
}

func (err Error) Error() string {
	return fmt.Sprintf("Received status code (%d) instead of 200 for a call to %s: %s", err.HTTPStatusCode, err.URL, err.error.Error())
}

//RateLimitedError struct to represent a Rate Limit has been hit for the given account
type RateLimitedError Error

func (err RateLimitedError) Error() string {
	return fmt.Sprintf("Your account has hit or exceded the current request/sec rate limitation")
}

//NewError constructes a new Client Error
func NewError(httpStatusCode int, url string, error error) error {
	switch httpStatusCode {
	case 429:
		return RateLimitedError{
			error:          error,
			HTTPStatusCode: httpStatusCode,
			URL:            url,
		}
	}

	return Error{
		error:          error,
		HTTPStatusCode: httpStatusCode,
		URL:            url,
	}
}
