package main

import (
	"fmt"
	"log"
	"io"
	"net/http"
	"encoding/json"
	"encoding/xml"
	//	"golang.org/x/net/html/charset"
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

func getCbrCurrencyXMLBody(curr *currencyAPIInput) ([]byte, error) {
	// need to add date
	endpoint := "http://www.cbr.ru/scripts/XML_daily.asp?date_req="

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		log.Printf("failed to create GET request for CBR %v", err)
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("failed to get http response: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("failed to read CBR response body: %v")
		return nil, err
	}

	return body, nil
}

func getCbrCurrencyValue(curr *currencyAPIInput) string {
	body, err := getCbrCurrencyXMLBody(curr)
	if err != nil {
		log.Printf("failed to get CBR currency XML body: %v\n")
		return ""
	}

	fmt.Printf("%s", body);

	return ""
}

func appCurrencyInfoHandler(w http.ResponseWriter, req *http.Request) {
	var input currencyAPIInput
	err := json.NewDecoder(req.Body).Decode(&input)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	value := getCbrCurrencyValue(&input)
	fmt.Printf("value %v\n", value)
}


func main() {
	http.HandleFunc("/info", appInfoHandler)
	http.HandleFunc("/info/currency", appCurrencyInfoHandler)
	http.ListenAndServe(":8080", nil)
}
