package handlers

import (
	"github.com/gin-gonic/gin"
	. "github.com/naskavinda/webpageanalyzer/models"
	"io"
	"net/http"
	"strings"
)

var HTTPGet = http.Get

func WebPageAnalyzerHandler(c *gin.Context) {
	var request PageAnalysisRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format or missing webpageUrl",
		})
		return
	}

	request.WebpageUrl = strings.TrimSpace(request.WebpageUrl)
	if request.WebpageUrl == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "URL cannot be empty",
		})
		return
	}

	resp, err := HTTPGet(request.WebpageUrl)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch the webpage: " + err.Error(),
		})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to fetch the webpage, status code: " + resp.Status,
		})
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to read the webpage content: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"url":     request.WebpageUrl,
		"content": string(body),
	})
}
