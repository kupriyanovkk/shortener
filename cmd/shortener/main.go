package main

import (
	"fmt"
	"net/http"

	"github.com/kupriyanovkk/shortener/internal/app"
)

var buildVersion string = "N/A"
var buildDate string = "N/A"
var buildCommit string = "N/A"

// printBuildInfo prints the build information.
func printBuildInfo() {
	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)
}

func init() {
	printBuildInfo()
}

func main() {
	go http.ListenAndServe(":9900", nil)

	app.Start()
}
