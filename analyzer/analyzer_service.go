package analyzer

import "github.com/naskavinda/webpageanalyzer/models"

type AnalyzerService interface {
	Analyze(url string) (models.PageAnalysisResponse, error)
}
