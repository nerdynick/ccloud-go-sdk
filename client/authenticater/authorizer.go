package authenticater

import "net/http"

type Authenticater interface {
	Authenticate(*http.Request)
}
