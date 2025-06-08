package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var originalHTTPGet = HTTPGet

func TestWebPageAnalyzerHandler_ValidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	req := newTestRequest(`{"webpageUrl": "  https://example.com  "}`)

	response := mockHTTPGetSuccess()
	cleanUp := setupMockHTTP(response, nil)
	defer cleanUp()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	WebPageAnalyzerHandler(c)

	assert.Equal(t, http.StatusOK, w.Code)

	resp := decodeJSONResponse(t, w.Body)
	assert.Equal(t, "https://example.com", resp["url"])
	assert.Equal(t, "<html><body>Test Page</body></html>", resp["content"])
}

func TestWebPageAnalyzerHandler_EmptyURL(t *testing.T) {
	gin.SetMode(gin.TestMode)

	req := newTestRequest(`{"webpageUrl": "  "}`)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	WebPageAnalyzerHandler(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	resp := decodeJSONResponse(t, w.Body)
	assert.Equal(t, "URL cannot be empty", resp["error"])
}

func TestWebPageAnalyzerHandler_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	req := newTestRequest(`{"webpageUrl": }`)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	WebPageAnalyzerHandler(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	resp := decodeJSONResponse(t, w.Body)
	assert.Equal(t, "Invalid request format or missing webpageUrl", resp["error"])
}

func TestWebPageAnalyzerHandler_InvalidURL(t *testing.T) {
	gin.SetMode(gin.TestMode)

	req := newTestRequest(`{"webpageUrl": "invalid-url"}`)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	WebPageAnalyzerHandler(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	resp := decodeJSONResponse(t, w.Body)
	assert.Equal(t, "Invalid URL format", resp["error"])
}

func decodeJSONResponse(t *testing.T, body *bytes.Buffer) map[string]string {
	t.Helper()
	var data map[string]string
	err := json.Unmarshal(body.Bytes(), &data)
	if err != nil {
		t.Fatalf("Failed to decode JSON response: %v", err)
	}
	return data
}

func newTestRequest(jsonBody string) *http.Request {
	req, _ := http.NewRequest(http.MethodPost, "/analyze", bytes.NewBufferString(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	return req
}

func mockHTTPGetSuccess() *http.Response {
	return newHTTPResponse("<html><body>Test Page</body></html>", http.StatusOK)
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
