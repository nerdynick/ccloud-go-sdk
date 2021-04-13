package responses

import "strings"

//ErrorResponse when a none 200 HTTP Status is returned. This handles the JSON
type ErrorResponse struct {
	Errors []Error `json:"errors"`
}

func (err ErrorResponse) Error() string {
	codes := make([]string, len(err.Errors))
	for i, e := range err.Errors {
		codes[i] = e.Code
	}
	return "Received the following errors from the CCloud API: " + strings.Join(codes, ",")
}

//Error a given error
type Error struct {
	ID         string        `json:"id,omitempty"`
	Status     string        `json:"status,omitempty"`
	Code       string        `json:"code,omitempty"`
	Title      string        `json:"title,omitempty"`
	Detail     string        `json:"detail"`
	Resolution string        `json:"resolution,omitempty"`
	Source     []ErrorSource `json:"source,omitempty"`
}

//ErrorSource error source representation
type ErrorSource struct {
	Pointer   string `json:"pointer,omitempty"`
	Parameter string `json:"parameter,omitempty"`
}
