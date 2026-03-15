package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"app/api"
)

var appInfo api.InfoAPIOutput

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = ":8000"
	} else {
		port = fmt.Sprintf(":%s", port)
	}

	author := os.Getenv("AUTHOR")
	if author == "" {
		author = "a.dashchinsky"
	}

	version := os.Getenv("VERSION")
	if version == "" {
		version = "1.0.0"
	}

	appInfo.Author = author
	appInfo.Version = version
	appInfo.Service = "currency"

	http.HandleFunc("/info", api.AppInfoHandler(appInfo))
	http.HandleFunc("/info/currency", api.AppCurrencyInfoHandler())

	log.Fatal(http.ListenAndServe(port, nil))
}
