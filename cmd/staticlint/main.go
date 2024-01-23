package main

import (
	"github.com/kupriyanovkk/shortener/pkg/staticlint"
	"golang.org/x/tools/go/analysis/multichecker"
)

func main() {
	analyzers := staticlint.GetAnalyzers()

	multichecker.Main(analyzers...)
}
