package analyzer

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

var originalHTTPGet = HTTPGet
var originalHTTPClient = HTTPClient

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

func TestAnalyze_ShouldReturnSuccessResponse_WithHeading(t *testing.T) {
	gin.SetMode(gin.TestMode)

	pageUrl := "https://example.com/test-page"

	response := newHTTPResponse(validHTMLContentWithHeaders, http.StatusOK)
	cleanUp := setupMockHTTP(response, nil)
	defer cleanUp()

	analyze, err := Analyze(pageUrl)

	assert.NoError(t, err)
	assert.Equal(t, pageUrl, analyze.URL)
	assert.Equal(t, "HTML5", analyze.HTMLVersion)
	assert.Equal(t, "Example Page with Various Links", analyze.Title)
	assert.Equal(t, 1, analyze.HeadingCounts["h1"])
	assert.Equal(t, 4, analyze.HeadingCounts["h2"])
}

func TestAnalyze_ShouldReturnSuccessResponse_WithoutHeading(t *testing.T) {
	gin.SetMode(gin.TestMode)

	pageUrl := "https://example.com/test-page"

	response := newHTTPResponse(validHTMLContentWithoutHeaders, http.StatusOK)
	cleanUp := setupMockHTTP(response, nil)
	defer cleanUp()

	analyze, err := Analyze(pageUrl)

	assert.NoError(t, err)
	assert.Equal(t, pageUrl, analyze.URL)
	assert.Equal(t, "HTML5", analyze.HTMLVersion)
	assert.Equal(t, "Page Without Headers", analyze.Title)
	assert.Equal(t, 0, len(analyze.HeadingCounts))
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

func TestLinksAnalyzer_ShouldReturnInternalAndExternalLinkCount(t *testing.T) {

	gin.SetMode(gin.TestMode)

	pageUrl := "https://example.com/test-page"

	parsedURL, err := url.Parse(pageUrl)
	assert.NoError(t, err)

	response := newHTTPResponse(validHTMLContentWithHeaders, http.StatusOK)
	cleanUp := setupMockHTTP(response, nil)
	defer cleanUp()

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(validHTMLContentWithHeaders))
	assert.NoError(t, err)

	internalCount, externalCount, _ := linksAnalyzer(doc, parsedURL)

	assert.Equal(t, 5, internalCount) // 2 internal links
	assert.Equal(t, 2, externalCount) // 2 external links
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

const validHTMLContentWithHeaders = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <title>Example Page with Various Links</title>
</head>
<body>
    <a id="top"></a>
    <h1>Welcome to the Example Page</h1>

    <nav>
        <h2>Navigation</h2>
        <ul>
            <li><a href="#section1">Go to Section 1</a></li>
            <li><a href="#section2">Go to Section 2</a></li>
            <li><a href="https://www.example.com" target="_blank" rel="noopener noreferrer">Visit Example.com (external)</a></li>
            <li><a>Inaccessible Link (no href)</a></li>
        </ul>
    </nav>

    <section id="section1">
        <h2>Section 1</h2>
        <p>This is the first section. You can <a href="#section2">jump to Section 2</a> or <a href="#top">go back to the top</a>.</p>
    </section>

    <section id="section2">
        <h2>Section 2</h2>
        <p>This is the second section. Here is an <a href="https://www.openai.com" target="_blank" rel="noopener noreferrer">external link to OpenAI</a>.</p>
    </section>

    <footer>
        <h2>Footer</h2>
        <p>Return to <a href="#top">top of page</a>.</p>
    </footer>
</body>
</html>`

const validHTMLContentWithoutHeaders = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <title>Page Without Headers</title>
</head>
<body>
    <a id="top"></a>

    <nav>
        <ul>
            <li><a href="#section1">Go to Section 1</a></li>
            <li><a href="#section2">Go to Section 2</a></li>
            <li><a href="https://www.example.com" target="_blank" rel="noopener noreferrer">Visit Example.com (external)</a></li>
            <li><a>Inaccessible Link (no href)</a></li>
        </ul>
    </nav>

    <section id="section1">
        <p>This is the first section. You can <a href="#section2">jump to Section 2</a> or <a href="#top">go back to the top</a>.</p>
    </section>

    <section id="section2">
        <p>This is the second section. Here is an <a href="https://www.openai.com" target="_blank" rel="noopener noreferrer">external link to OpenAI</a>.</p>
    </section>

    <footer>
        <p>Return to <a href="#top">top of page</a>.</p>
    </footer>
</body>
</html>
`
