package main

import (
	"fmt"
	"net/http"
	"encoding/json"
)

type appInfo struct {
	Version string `json:"version"`
	Service string `json:"service"`
	Author string `json:"author"`
}

type currencyAPIInput struct {
	Currency string `json:"currency"`
	Date string `json:"date"`
}

func appInfoHandler(w http.ResponseWriter, req *http.Request) {
	info := appInfo{Version: "0.1.0", Service: "currency", Author: "a.dashchynski"}
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(info); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return;
	}
}

func appCurrencyInfoHandler(w http.ResponseWriter, req *http.Request) {
	var input currencyAPIInput
	err := json.NewDecoder(req.Body).Decode(&input)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, "You try to get currency of %s on %s.\n",
		input.Currency, input.Date)
}


func main() {
	http.HandleFunc("/info", appInfoHandler)
	http.HandleFunc("/info/currency", appCurrencyInfoHandler)
	http.ListenAndServe(":8080", nil)
}
