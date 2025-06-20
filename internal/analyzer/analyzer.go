package analyzer

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	. "github.com/naskavinda/webpageanalyzer/internal/model"
	"github.com/naskavinda/webpageanalyzer/internal/validator"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

var HTTPGet = http.Get
var HTTPClient = http.Client{}

type DefaultAnalyzerService struct{}

func (defaultAnalyzer DefaultAnalyzerService) Analyze(pageUrl string) (PageAnalysisResponse, error) {
	log.Printf("[DEBUG] Starting analysis for URL: %s", pageUrl)
	var isValidURL = false

	isValidURL = validator.IsValidURL(&pageUrl)

	if !isValidURL {
		log.Printf("[ERROR] Invalid URL format: %s", pageUrl)
		return PageAnalysisResponse{}, fmt.Errorf("invalid URL format")
	}

	resp, err := HTTPGet(pageUrl)
	if err != nil {
		log.Printf("[ERROR] Failed to fetch the webpage: %v", err)
		return PageAnalysisResponse{}, fmt.Errorf("failed to fetch the webpage")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("[ERROR] Non-200 status code for %s: %v", pageUrl, resp.Status)
		return PageAnalysisResponse{}, fmt.Errorf("failed to fetch the webpage, status code: %v", resp.Status)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Printf("[ERROR] Failed to read the webpage content for %s: %v", pageUrl, err)
		return PageAnalysisResponse{}, fmt.Errorf("failed to read the webpage content")
	}

	result := PageAnalysisResponse{
		URL:           pageUrl,
		HeadingCounts: make(map[string]int),
	}

	result.HTMLVersion = detectHTMLVersion(doc)
	log.Printf("[DEBUG] Detected HTML version for %s: %s", pageUrl, result.HTMLVersion)

	result.Title = doc.Find("title").Text()
	log.Printf("[DEBUG] Page title for %s: %s", pageUrl, result.Title)

	getHeadingCount(doc, result)

	parsedURL, err := getUrl(pageUrl, err)
	if err != nil {
		log.Printf("[ERROR] Failed to parse URL %s: %v", pageUrl, err)
		return PageAnalysisResponse{}, err
	}

	internalCount, externalCount, inaccessibleCount := linksAnalyzer(doc, parsedURL)

	result.InternalLinks = internalCount
	result.ExternalLinks = externalCount
	result.InaccessibleLinks = inaccessibleCount

	result.HasLoginForm = detectLoginForm(doc)

	log.Printf("[INFO] Analysis complete for %s", pageUrl)
	return result, nil
}

func getUrl(pageUrl string, err error) (*url.URL, error) {
	parsedURL, err := url.Parse(pageUrl)
	if err != nil {
		log.Printf("[ERROR] url.Parse failed for %s: %v", pageUrl, err)
		return nil, fmt.Errorf("given URL is invalid")
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
			log.Printf("[ERROR] Failed to parse link href: %s, error: %v", href, err)
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
			wg.Add(1)
			go func(link string) {

				defer wg.Done()

				if !isLinkAccessible(link) {
					log.Printf("[DEBUG] Link inaccessible: %s", link)
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
		log.Printf("[DEBUG] Link not accessible: %s, err: %v, status: %v", link, err, resp)
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
