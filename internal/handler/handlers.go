package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/naskavinda/webpageanalyzer/internal/analyzer"
	. "github.com/naskavinda/webpageanalyzer/internal/model"
)

type WebPageAnalyzer struct {
	Service analyzer.AnalyzerService
}

func (webPageAnalyzer *WebPageAnalyzer) WebPageAnalyzerHandler(c *gin.Context) {
	var request PageAnalysisRequest

	log.Println("[INFO] Received /analyzer request")

	if err := c.ShouldBindJSON(&request); err != nil {
		log.Printf("[ERROR] Invalid request format or missing webpageUrl: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format or missing webpageUrl",
		})
		return
	}
	response, err := webPageAnalyzer.Service.Analyze(request.WebpageUrl)
	if err != nil {
		log.Printf("[ERROR] Analysis failed for %s: %v", request.WebpageUrl, err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	log.Printf("[INFO] Analysis successful for %s", request.WebpageUrl)
	c.JSON(http.StatusOK, gin.H{
		"url":     request.WebpageUrl,
		"content": response,
	})
}
