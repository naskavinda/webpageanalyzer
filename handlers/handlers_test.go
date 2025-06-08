package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

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
