package main

import (
	"fmt"
	"log"
	"io"
	"net/http"
	"encoding/json"
	"encoding/xml"
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

type Valute struct {
	ID       string `xml:"ID,attr"`
	CharCode string `xml:"CharCode"`
	Value    string `xml:"Value"`
	Nominal  string `xml:"Nominal"`
}

type ValCurs struct {
	XMLName xml.Name `xml:"ValCurs"`
	Valutes [] Valute `xml:"Valute"`
}

func findCurrency(curs ValCurs, valName string) *Valute {
	for _, v := range curs.Valutes {
		fmt.Printf("%s\n", v.CharCode)
		if v.CharCode == valName {
			return &v
		}
	}

	return nil
}

func getCurrencyValue(curr *currencyAPIInput) string {
	return ""
}

func appCurrencyInfoHandler(w http.ResponseWriter, req *http.Request) {
	var input currencyAPIInput
	err := json.NewDecoder(req.Body).Decode(&input)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	value := getCurrencyValue(&input)
}


func main() {
	http.HandleFunc("/info", appInfoHandler)
	http.HandleFunc("/info/currency", appCurrencyInfoHandler)
	http.ListenAndServe(":8080", nil)
}
