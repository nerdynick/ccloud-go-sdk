package client

const (
	//DefaultUserAgent is the default user agent to send
	DefaultUserAgent string = "ccloud-metrics-sdk/go"
)

//Context is the Contextual set of configs for the HTTP Client making the calls to the Metrics API
type Context struct {
	APIKey      string
	APISecret   string
	UserAgent   string
	HTTPHeaders map[string]string
}

//NewContext creates a new instance of the HTTPContext loaded with the defaults where possible
func NewContext(apiKey string, apiSecret string) Context {
	return Context{
		UserAgent:   DefaultUserAgent,
		HTTPHeaders: nil,
		APIKey:      apiKey,
		APISecret:   apiSecret,
	}
}
