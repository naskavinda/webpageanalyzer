package analyzer

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	. "github.com/naskavinda/webpageanalyzer/models"
	"github.com/naskavinda/webpageanalyzer/validaters"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

var HTTPGet = http.Get
var HTTPClient = http.Client{}

type DefaultAnalyzerService struct{}

func (defaultAnalyzer DefaultAnalyzerService) Analyze(pageUrl string) (PageAnalysisResponse, error) {

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

	parsedURL, err := getUrl(pageUrl, err)
	if err != nil {
		return PageAnalysisResponse{}, err
	}

	internalCount, externalCount, inaccessibleCount := linksAnalyzer(doc, parsedURL)

	result.InternalLinks = internalCount
	result.ExternalLinks = externalCount
	result.InaccessibleLinks = inaccessibleCount

	result.HasLoginForm = detectLoginForm(doc)

	return result, nil
}

func getUrl(pageUrl string, err error) (*url.URL, error) {
	parsedURL, err := url.Parse(pageUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %v", err.Error())
	}
	return parsedURL, nil
}

func linksAnalyzer(doc *goquery.Document, baseUrl *url.URL) (int, int, int) {

	var internalCount, externalCount, inaccessibleCount int
	var wg sync.WaitGroup
	var mu sync.Mutex

	links := doc.Find("a[href]")

	links.Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")

		if !exists || href == "" || href == "#" {
			return
		}

		linkUrl, err := url.Parse(href)

		if err != nil {
			mu.Lock()
			inaccessibleCount++
			mu.Unlock()
			return
		}

		if !linkUrl.IsAbs() {
			linkUrl = baseUrl.ResolveReference(linkUrl)
		}

		if linkUrl.Host == baseUrl.Host {
			mu.Lock()
			internalCount++
			mu.Unlock()
		} else {
			mu.Lock()
			externalCount++
			mu.Unlock()
			go func(link string) {
				wg.Add(1)
				defer wg.Done()

				if !isLinkAccessible(link) {
					mu.Lock()
					inaccessibleCount++
					mu.Unlock()
				}

			}(linkUrl.String())
		}

	})
	wg.Wait()
	return internalCount, externalCount, inaccessibleCount
}

func isLinkAccessible(link string) bool {
	resp, err := HTTPClient.Head(link)
	if err != nil || resp.StatusCode >= 400 {
		return false
	}
	if resp != nil {
		resp.Body.Close()
	}
	return true
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

func detectLoginForm(doc *goquery.Document) bool {

	loginTexts := []string{"login", "log in", "sign in", "sign up"}
	hasLoginText := false

	formHTML, _ := doc.Html()
	formText := strings.ToLower(formHTML)
	hasPasswordField := strings.Contains(formText, "type='password'") || strings.Contains(formText, "type=\"password\"")

	for _, text := range loginTexts {
		if strings.Contains(formText, text) {
			hasLoginText = true
			break
		}
	}

	action, _ := doc.Attr("action")
	id, _ := doc.Attr("id")

	for _, text := range loginTexts {
		if strings.Contains(strings.ToLower(action), text) || strings.Contains(strings.ToLower(id), text) {
			hasLoginText = true
			break
		}
	}

	return hasPasswordField || hasLoginText
}
