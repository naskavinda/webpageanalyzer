package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/naskavinda/webpageanalyzer/internal/analyzer"
	. "github.com/naskavinda/webpageanalyzer/internal/handler"
	"log"
	"time"
)

func main() {
	log.Println("[INFO] Starting Web Page Analyzer server...")
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	w := WebPageAnalyzer{
		Service: analyzer.DefaultAnalyzerService{},
	}
	log.Println("[INFO] Registering /analyzer endpoint")
	r.POST("/analyzer", w.WebPageAnalyzerHandler)

	log.Println("[INFO] Server is running on :8080")
	r.Run()
}
