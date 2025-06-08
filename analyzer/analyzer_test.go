package analyzer

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"strings"
	"testing"
)

var originalHTTPGet = HTTPGet

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

func TestAnalyze_ShouldGiveValidHTMLVersion(t *testing.T) {
	gin.SetMode(gin.TestMode)

	pageUrl := "https://example.com/test-page"

	response := mockHTTPGetSuccess()
	cleanUp := setupMockHTTP(response, nil)

	analyze, err := Analyze(pageUrl)

	defer cleanUp()

	assert.NoError(t, err)
	assert.Equal(t, pageUrl, analyze.URL)
	assert.Equal(t, "HTML5", analyze.HTMLVersion)
}

func TestAnalyze_ShouldReturnErrorOnHTTPFailure(t *testing.T) {
	gin.SetMode(gin.TestMode)

	pageUrl := "https://example.com/test-page"

	mockError := http.ErrHandlerTimeout
	cleanUp := setupMockHTTP(nil, mockError)
	defer cleanUp()

	_, err := Analyze(pageUrl)

	assert.Error(t, err)
	assert.Equal(t, "failed to fetch the webpage: http: Handler timeout", err.Error())
}

func TestDetectHTMLVersion(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		expected string
	}{
		{
			name:     "HTML5",
			html:     "<!DOCTYPE html><html><head></head><body></body></html>",
			expected: "HTML5",
		},
		{
			name:     "HTML 4.01",
			html:     "<!DOCTYPE HTML PUBLIC \"-//W3C//DTD HTML 4.01 Transitional//EN\"><html><body></body></html>",
			expected: "HTML 4.01",
		},
		{
			name:     "XHTML",
			html:     "<!DOCTYPE html PUBLIC \"-//W3C//DTD XHTML 1.0 Transitional//EN\"><html><body></body></html>",
			expected: "XHTML 1.0",
		},
		{
			name:     "XHTML",
			html:     "<!DOCTYPE html PUBLIC \"-//W3C//DTD XHTML 1.1 Transitional//EN\"><html><body></body></html>",
			expected: "XHTML 1.1",
		},
		{
			name:     "Unknown",
			html:     "<!DOCTYPE something-custom><html><body></body></html>",
			expected: "Unknown",
		},
		{
			name:     "No DOCTYPE",
			html:     "<html><body>No doctype</body></html>",
			expected: "Unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			isHtmlVersionCorrect(t, tt.html, tt.expected)
		})
	}
}

func isHtmlVersionCorrect(t *testing.T, htmlContent string, expectedVersion string) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	assert.NoError(t, err)

	version := detectHTMLVersion(doc)
	assert.Equal(t, expectedVersion, version)
}

func mockHTTPGetSuccess() *http.Response {
	return newHTTPResponse("<!DOCTYPE html><body>Test Page</body></html>", http.StatusOK)
}

func newHTTPResponse(response string, statusCode int) *http.Response {
	body := io.NopCloser(strings.NewReader(response))
	return &http.Response{
		StatusCode: statusCode,
		Body:       body,
	}
}

func setupMockHTTP(mockResponse *http.Response, mockError error) func() {
	HTTPGet = func(url string) (*http.Response, error) {
		return mockResponse, mockError
	}

	return func() {
		HTTPGet = originalHTTPGet
	}
}
