package client

const (
	//DefaultUserAgent is the default user agent to send
	DefaultUserAgent string = "ccloud-metrics-sdk/go"
)

//Context is the Contextual set of configs for the HTTP Client making the calls to the Metrics API
type Context struct {
	UserAgent   string
	HTTPHeaders map[string]string
	BaseURL     string
}

//NewContext creates a new instance of the HTTPContext loaded with the defaults where possible
func NewContext(baseURL string) Context {
	return Context{
		UserAgent:   DefaultUserAgent,
		HTTPHeaders: nil,
		BaseURL:     baseURL,
	}
}
