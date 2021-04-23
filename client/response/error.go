package response

import (
	"fmt"
	"strings"
)

//ErrorResponse when a none 200 HTTP Status is returned. This handles the JSON
type ErrorResponse struct {
	Errors []Error `json:"errors,omitempty"`
	Err    Error   `json:"error,omitempty"`
}

func (err ErrorResponse) Error() string {
	if len(err.Errors) > 1 {
		codes := make([]string, len(err.Errors))
		for i, e := range err.Errors {
			codes[i] = e.Error()
		}
		return fmt.Sprintf("CCloud API Errors [`%s`]", strings.Join(codes, "`,`"))
	} else if len(err.Errors) == 1 {
		return err.Errors[0].Error()
	} else {
		return err.Err.Error()
	}

}

//Error a given error
type Error struct {
	ID         string        `json:"id,omitempty"`
	Status     string        `json:"status,omitempty"`
	Code       string        `json:"code,omitempty"`
	Message    string        `json:"message,omitempty"`
	Title      string        `json:"title,omitempty"`
	Detail     string        `json:"detail,omitempty"`
	Resolution string        `json:"resolution,omitempty"`
	Source     []ErrorSource `json:"source,omitempty"`
}

func (err Error) Error() string {
	if err.Message != "" {
		return fmt.Sprintf("APIError(%s - %s)", err.Code, err.Message)
	} else {
		return fmt.Sprintf("APIError(%s/%s - %s)", err.Code, err.Status, err.Detail)
	}
}

//ErrorSource error source representation
type ErrorSource struct {
	Pointer   string `json:"pointer,omitempty"`
	Parameter string `json:"parameter,omitempty"`
}
