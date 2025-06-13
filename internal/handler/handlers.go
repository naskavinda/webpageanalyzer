package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/naskavinda/webpageanalyzer/internal/analyzer"
	. "github.com/naskavinda/webpageanalyzer/internal/model"
	"net/http"
)

type WebPageAnalyzer struct {
	Service analyzer.AnalyzerService
}

func (webPageAnalyzer *WebPageAnalyzer) WebPageAnalyzerHandler(c *gin.Context) {
	var request PageAnalysisRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format or missing webpageUrl",
		})
		return
	}
	response, err := webPageAnalyzer.Service.Analyze(request.WebpageUrl)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"url":     request.WebpageUrl,
		"content": response,
	})
}
