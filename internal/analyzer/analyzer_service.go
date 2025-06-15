package analyzer

import (
	"github.com/naskavinda/webpageanalyzer/internal/model"
)

type Service interface {
	Analyze(url string) (model.PageAnalysisResponse, error)
}
