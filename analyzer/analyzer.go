package analyzer

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	. "github.com/naskavinda/webpageanalyzer/models"
	"github.com/naskavinda/webpageanalyzer/validaters"
	"net/http"
	"strings"
)

var HTTPGet = http.Get

func Analyze(pageUrl string) (PageAnalysisResponse, error) {

	var isValidURL = false

	isValidURL = validaters.IsValidURL(&pageUrl)

	if !isValidURL {
		return PageAnalysisResponse{}, fmt.Errorf("invalid URL format")
	}

	resp, err := HTTPGet(pageUrl)
	if err != nil {
		return PageAnalysisResponse{}, fmt.Errorf("failed to fetch the webpage: %v", err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return PageAnalysisResponse{}, fmt.Errorf("failed to fetch the webpage, status code: %v", resp.Status)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return PageAnalysisResponse{}, fmt.Errorf("failed to read the webpage content: %v", err.Error())
	}

	result := PageAnalysisResponse{
		URL:           pageUrl,
		HeadingCounts: make(map[string]int),
	}

	result.HTMLVersion = detectHTMLVersion(doc)

	result.Title = doc.Find("title").Text()

	getHeadingCount(doc, result)

	return result, nil
}

func getHeadingCount(doc *goquery.Document, result PageAnalysisResponse) {
	for i := 0; i < 7; i++ {
		selector := fmt.Sprintf("h%d", i)
		count := doc.Find(selector).Length()
		if count > 0 {
			result.HeadingCounts[selector] = count
		}
	}
}

func detectHTMLVersion(doc *goquery.Document) string {
	html, err := doc.Html()
	if err != nil {
		return "Unknown"
	}

	lowerCaseHTML := strings.ToLower(html)
	if strings.Contains(lowerCaseHTML, "<!doctype html>") {
		return "HTML5"
	}

	if strings.Contains(lowerCaseHTML, "html 4.01") {
		return "HTML 4.01"
	}

	if strings.Contains(lowerCaseHTML, "xhtml 1.0") {
		return "XHTML 1.0"
	}
	if strings.Contains(lowerCaseHTML, "xhtml 1.1") {
		return "XHTML 1.1"
	}

	return "Unknown"
}
