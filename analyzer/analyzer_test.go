package analyzer

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAnalyze_ValidPageURL(t *testing.T) {
	gin.SetMode(gin.TestMode)

	pageUrl := "https://google.com"
	analyze, err := Analyze(pageUrl)

	assert.NoError(t, err)
	assert.Equal(t, pageUrl, analyze.URL)
}

func TestAnalyze_InvalidPageURL(t *testing.T) {
	gin.SetMode(gin.TestMode)

	pageUrl := "invalid-url"
	_, err := Analyze(pageUrl)

	assert.Error(t, err)
	assert.Equal(t, "invalid URL format", err.Error())
}

func TestAnalyze_PageNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	pageUrl := "https://example.com/404"
	_, err := Analyze(pageUrl)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to fetch the webpage, status code: 404")
}

func TestAnalyze_EmptyPageURL(t *testing.T) {
	gin.SetMode(gin.TestMode)

	pageUrl := ""
	_, err := Analyze(pageUrl)

	assert.Error(t, err)
	assert.Equal(t, "invalid URL format", err.Error())
}

func TestAnalyze_URLWithLeadingSpaces(t *testing.T) {
	gin.SetMode(gin.TestMode)

	pageUrl := "   https://google.com"
	analyze, err := Analyze(pageUrl)

	assert.NoError(t, err)
	assert.Equal(t, "https://google.com", analyze.URL)
}

func TestAnalyze_URLWithTrailingSpaces(t *testing.T) {
	gin.SetMode(gin.TestMode)

	pageUrl := "https://google.com   "
	analyze, err := Analyze(pageUrl)

	assert.NoError(t, err)
	assert.Equal(t, "https://google.com", analyze.URL)
}
