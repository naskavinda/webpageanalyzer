package analyzer

import (
	"github.com/naskavinda/webpageanalyzer/internal/model"
)

type AnalyzerService interface {
	Analyze(url string) (model.PageAnalysisResponse, error)
}
