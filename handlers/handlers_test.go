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

func TestWebPageAnalyzerHandler_ValidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	HTTPGet = func(url string) (*http.Response, error) {
		body := io.NopCloser(strings.NewReader("<html><body>Test Page</body></html>"))
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       body,
		}, nil
	}

	jsonBody := `{"webpageUrl": "  https://example.com  "}`
	req, _ := http.NewRequest(http.MethodPost, "/analyze", bytes.NewBufferString(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	WebPageAnalyzerHandler(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]string
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Contains(t, "https://example.com", resp["url"])
	assert.Contains(t, "<html><body>Test Page</body></html>", resp["content"])
}
