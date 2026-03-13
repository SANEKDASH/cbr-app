package cbr

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"golang.org/x/net/html/charset"
)

type CbrValute struct {
	ID       string `xml:"ID,attr"`
	CharCode string `xml:"CharCode"`
	Value    string `xml:"Value"`
	Nominal  string `xml:"Nominal"`
}

type CbrValCurs struct {
	XMLName xml.Name    `xml:"ValCurs"`
	Valutes []CbrValute `xml:"Valute"`
}

const cbrDefaultCurrencyEndpoint = "http://www.cbr.ru/scripts/XML_daily.asp"

func GetCbrEndpoint(date string) (string, error) {
	if date == "" {
		return cbrDefaultCurrencyEndpoint, nil
	}

	dateParts := strings.Split(date, "-")
	if len(dateParts) != 3 {
		return "", fmt.Errorf("invalid date format: %s", date)
	}

	year := dateParts[0]
	month := dateParts[1]
	day := dateParts[2]

	cbrDate := fmt.Sprintf("%s/%s/%s", day, month, year)
	return fmt.Sprintf("%s?date_req=%s",
		cbrDefaultCurrencyEndpoint, cbrDate), nil
}

func GetCbrCurrencyXMLBody(date string) ([]byte, error) {
	endpoint, err := GetCbrEndpoint(date)
	if err != nil {
		log.Printf("failed to create endpoint for CBR: %v", err)
		return nil, err
	}

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
		log.Printf("failed to read CBR response body: %v", err)
		return nil, err
	}

	return body, nil
}

func GetCbrCurrencyValues(date string) (CbrValCurs, error) {
	body, err := GetCbrCurrencyXMLBody(date)
	if err != nil {
		log.Printf("failed to get CBR currency XML body: %v", err)
		return CbrValCurs{}, err
	}

	var curs CbrValCurs
	decoder := xml.NewDecoder(bytes.NewReader(body))
	decoder.CharsetReader = charset.NewReaderLabel

	err = decoder.Decode(&curs)
	if err != nil {
		log.Printf("failed to decode XML: %v", err)
		return CbrValCurs{}, err
	}

	return curs, nil
}

func FindCbrValute(curs CbrValCurs, valName string) *CbrValute {
	for _, v := range curs.Valutes {
		if v.CharCode == valName {
			return &v
		}
	}

	return nil
}

func ParseCbrValue(value string) (float64, error) {
	return strconv.ParseFloat(strings.ReplaceAll(value, ",", "."), 64)
}
