package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	. "github.com/naskavinda/webpageanalyzer/handlers"
)

func main() {
	fmt.Println("Welcome to the Web Page Analyzer!")
	r := gin.Default()

	r.POST("/analyzer", WebPageAnalyzerHandler)

	r.Run()
}
