package validaters

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsValidURL_URLEmpty(c *testing.T) {
	uri := " "
	isValid := IsValidURL(&uri)

	assert.Equal(c, isValid, false)
}

func TestIsValidURL_URLSchemaEmpty(c *testing.T) {
	uri := "example.com"
	isValid := IsValidURL(&uri)

	assert.Equal(c, isValid, false)
}

func TestIsValidURL_URLHTTPValid(c *testing.T) {
	uri := "http://example.com"
	isValid := IsValidURL(&uri)

	assert.Equal(c, isValid, true)
}

func TestIsValidURL_URLHTTPSValid(c *testing.T) {
	uri := "https://example.com"
	isValid := IsValidURL(&uri)

	assert.Equal(c, isValid, true)
}
