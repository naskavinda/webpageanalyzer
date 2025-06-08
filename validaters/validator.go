package validaters

import (
	"net/url"
	"strings"
)

func IsValidURL(uri string) (string, bool) {
	uri = strings.TrimSpace(uri)
	parsedURL, err := url.ParseRequestURI(uri)
	return uri, err == nil && (parsedURL.Scheme == "http" || parsedURL.Scheme == "https")
}
