package models

type PageAnalysisRequest struct {
	WebpageUrl string `json:"webpageUrl" binding:"required"`
}

type PageAnalysisResponse struct {
	URL               string
	HTMLVersion       string
	Title             string
	HeadingCounts     map[string]int
	InternalLinks     int
	ExternalLinks     int
	InaccessibleLinks int
	HasLoginForm      bool
	Error             string
}
