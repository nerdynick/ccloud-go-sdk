package client

import (
	"fmt"
	"strings"
)

type APIPath string

func (p APIPath) Format(client Client, apiVersion int8) string {
	return strings.Join([]string{
		client.Context.BaseURL,
		"v" + fmt.Sprint(apiVersion),
		string(p),
	}, "/")
}
