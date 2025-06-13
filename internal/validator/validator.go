package validator

import (
	"net/url"
	"strings"
)

func IsValidURL(uri *string) bool {
	*uri = strings.TrimSpace(*uri)
	parsedURL, err := url.ParseRequestURI(*uri)
	return err == nil && (parsedURL.Scheme == "http" || parsedURL.Scheme == "https")
}
