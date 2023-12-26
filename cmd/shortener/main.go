package main

import (
	"net/http"

	"github.com/kupriyanovkk/shortener/internal/app"
)

func main() {
	go http.ListenAndServe(":9900", nil)

	app.Start()
}
