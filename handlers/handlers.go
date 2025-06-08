package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/naskavinda/webpageanalyzer/analyzer"
	. "github.com/naskavinda/webpageanalyzer/models"
	"net/http"
)

func WebPageAnalyzerHandler(c *gin.Context) {
	var request PageAnalysisRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format or missing webpageUrl",
		})
		return
	}
	response, err := analyzer.Analyze(request.WebpageUrl)
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
