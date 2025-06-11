package validaters

import (
	"github.com/gin-gonic/gin"
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

func TestAnalyze_URLWithLeadingSpaces(t *testing.T) {
	gin.SetMode(gin.TestMode)

	pageUrl := "   https://google.com"
	isValid := IsValidURL(&pageUrl)

	assert.Equal(t, isValid, true)
	assert.Equal(t, "https://google.com", pageUrl)
}

func TestAnalyze_URLWithTrailingSpaces(t *testing.T) {
	gin.SetMode(gin.TestMode)

	pageUrl := "https://google.com   "
	isValid := IsValidURL(&pageUrl)

	assert.Equal(t, isValid, true)
	assert.Equal(t, "https://google.com", pageUrl)
}
