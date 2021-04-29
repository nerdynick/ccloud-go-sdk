package authenticater

import (
	"net/http"

	"github.com/nerdynick/ccloud-go-sdk/client"
)

type APIKeyAuthenticater struct {
	APIKey    string
	APISecret client.SecurePassword
}

func (a *APIKeyAuthenticater) UpdateKey(key string) {
	a.APIKey = key
}

func (a *APIKeyAuthenticater) UpdateSecret(key string) {
	a.APISecret = client.SecurePassword(key)
}

func (a APIKeyAuthenticater) Authorize(req *http.Request) {
	req.SetBasicAuth(a.APIKey, a.APISecret.Value())
}

func NewAPIKeyAuth(apiKey string, apiSecret string) APIKeyAuthenticater {
	return APIKeyAuthenticater{
		APIKey:    apiKey,
		APISecret: client.SecurePassword(apiSecret),
	}
}
