package models

type PageAnalysisRequest struct {
	WebpageUrl string `json:"webpageUrl" binding:"required"`
}
