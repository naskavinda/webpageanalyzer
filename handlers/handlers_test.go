package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/naskavinda/webpageanalyzer/models"
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
	mockService := MockAnalyzerService{
		AnalyzeFunc: func(url string) (models.PageAnalysisResponse, error) {
			return models.PageAnalysisResponse{}, fmt.Errorf("Invalid request format or missing webpageUrl")
		},
	}
	var webPageAnalyzer = WebPageAnalyzer{Service: mockService}
	webPageAnalyzer.WebPageAnalyzerHandler(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	resp := decodeJSONResponse(t, w.Body)
	assert.Equal(t, "Invalid request format or missing webpageUrl", resp["error"])
}

func TestWebPageAnalyzerHandler_InvalidWebPageURL(t *testing.T) {
	gin.SetMode(gin.TestMode)

	req := newTestRequest(`{"webpageUrl":"https://example.com" }`)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	mockService := MockAnalyzerService{
		AnalyzeFunc: func(url string) (models.PageAnalysisResponse, error) {
			return models.PageAnalysisResponse{}, fmt.Errorf("invalid URL format")
		},
	}
	var webPageAnalyzer = WebPageAnalyzer{Service: mockService}
	webPageAnalyzer.WebPageAnalyzerHandler(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	resp := decodeJSONResponse(t, w.Body)
	assert.Equal(t, "invalid URL format", resp["error"])
}

func TestWebPageAnalyzerHandler_validJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	req := newTestRequest(`{"webpageUrl":"https://example.com" }`)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	mockService := MockAnalyzerService{
		AnalyzeFunc: func(url string) (models.PageAnalysisResponse, error) {
			return models.PageAnalysisResponse{
				URL:         url,
				HTMLVersion: "HTML5",
				Title:       "Sample Title",
			}, nil
		},
	}
	var webPageAnalyzer = WebPageAnalyzer{Service: mockService}
	webPageAnalyzer.WebPageAnalyzerHandler(c)

	assert.Equal(t, http.StatusOK, w.Code)

	resp := decodePageAnalysisResponse(t, w.Body)
	assert.Equal(t, "https://example.com", resp.URL)
}

func decodePageAnalysisResponse(t *testing.T, body *bytes.Buffer) models.PageAnalysisResponse {
	t.Helper()
	var data models.PageAnalysisResponse
	err := json.Unmarshal(body.Bytes(), &data)
	if err != nil {
		t.Fatalf("Failed to decode JSON response: %v", err)
	}
	return data
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

type MockAnalyzerService struct {
	AnalyzeFunc func(url string) (models.PageAnalysisResponse, error)
}

func (s MockAnalyzerService) Analyze(url string) (models.PageAnalysisResponse, error) {
	if s.AnalyzeFunc != nil {
		return s.AnalyzeFunc(url)
	}
	return models.PageAnalysisResponse{}, nil
}
