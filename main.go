package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/naskavinda/webpageanalyzer/analyzer"
	. "github.com/naskavinda/webpageanalyzer/handlers"
	"time"
)

func main() {
	fmt.Println("Welcome to the Web Page Analyzer!")
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
	r.POST("/analyzer", w.WebPageAnalyzerHandler)

	r.Run()
}
