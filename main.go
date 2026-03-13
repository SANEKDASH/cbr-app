package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"strings"
	"strconv"
	"io"
	"fmt"
	"log"
	"net/http"
	"golang.org/x/net/html/charset"
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

type currencyAPIOutput struct {
	Service string `json:"service"`
	Data map[string]float64 `json:"data"`

}

func NewCurrencyAPIOutput() *currencyAPIOutput {
	return &currencyAPIOutput{
		Service: "currency",
		Data:    make(map[string]float64),
	}
}

func appInfoHandler(w http.ResponseWriter, req *http.Request) {
	info := appInfo{Version: "0.1.0", Service: "currency", Author: "a.dashchynski"}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(info); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return;
	}
}

type CbrValute struct {
	ID       string `xml:"ID,attr"`
	CharCode string `xml:"CharCode"`
	Value    string `xml:"Value"`
	Nominal  string `xml:"Nominal"`
}

type CbrValCurs struct {
	XMLName xml.Name `xml:"ValCurs"`
	Valutes [] CbrValute `xml:"Valute"`
}

func findCbrValute(curs CbrValCurs, valName string) *CbrValute {
	for _, v := range curs.Valutes {
		if v.CharCode == valName {
			return &v
		}
	}

	return nil
}

func getCbrCurrencyXMLBody(curr currencyAPIInput) ([]byte, error) {
	endpoint := fmt.Sprintf("http://www.cbr.ru/scripts/XML_daily.asp?date_req=%s",
		strings.ReplaceAll(curr.Date, "-", "/"))

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

func addValuteToCurrencyAPIOutput(out *currencyAPIOutput, v CbrValute) error {
	val, err := strconv.ParseFloat(strings.ReplaceAll(v.Value, ",", "."), 64)
	if err != nil {
		log.Printf("failed to get value of currency: %v", err)
		return err
	}

	out.Data[v.CharCode] = val

	return nil
}

func fillCurrencyAPIOutput(out *currencyAPIOutput, curs CbrValCurs) error {
	for _, v := range curs.Valutes {
		err := addValuteToCurrencyAPIOutput(out, v)
		if err != nil {
			return err
		}
	}
	return nil
}

func getCbrCurrencyValues(curr currencyAPIInput) (CbrValCurs, error) {
	body, err := getCbrCurrencyXMLBody(curr)
	if err != nil {
		log.Printf("failed to get CBR currency XML body: %v\n")
		return CbrValCurs{}, err
	}

	var curs CbrValCurs
	decoder := xml.NewDecoder(bytes.NewReader(body))
	decoder.CharsetReader = charset.NewReaderLabel

	err = decoder.Decode(&curs)
	if err != nil {
		log.Printf("failed to decode XML: %v\n", err)
		return CbrValCurs{}, err
	}

	return curs, nil
}

func appCurrencyInfoHandler(w http.ResponseWriter, req *http.Request) {
	var input currencyAPIInput
	err := json.NewDecoder(req.Body).Decode(&input)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	values, err := getCbrCurrencyValues(input)
	if err != nil {
		http.Error(w, "failed to get currency values", http.StatusNoContent)
		return
	}

	output := NewCurrencyAPIOutput()
	if input.Currency == "" {
		fillCurrencyAPIOutput(output, values)
	} else {
		valute := findCbrValute(values, input.Currency)
		if valute == nil {
			http.Error(w, "not found", http.StatusNotFound)
			return;
		}

		addValuteToCurrencyAPIOutput(output, *valute)
	}
	if err := json.NewEncoder(w).Encode(output); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return;
	}
}

func main() {
	http.HandleFunc("/info", appInfoHandler)
	http.HandleFunc("/info/currency", appCurrencyInfoHandler)
	http.ListenAndServe(":8080", nil)
}
