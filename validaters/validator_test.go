package validaters

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsValidURL_URLEmpty(c *testing.T) {
	isValid := IsValidURL(" ")

	assert.Equal(c, isValid, false)
}

func TestIsValidURL_URLSchemaEmpty(c *testing.T) {
	isValid := IsValidURL("example.com")

	assert.Equal(c, isValid, false)
}

func TestIsValidURL_URLHTTPValid(c *testing.T) {
	isValid := IsValidURL("http://example.com")

	assert.Equal(c, isValid, true)
}

func TestIsValidURL_URLHTTPSValid(c *testing.T) {
	isValid := IsValidURL("https://example.com")

	assert.Equal(c, isValid, true)
}
