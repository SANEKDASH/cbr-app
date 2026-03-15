package api

import (
	"encoding/json"
	"log"
	"net/http"

	"app/cbr"
)

type InfoAPIOutput struct {
	Version string `json:"version"`
	Service string `json:"service"`
	Author  string `json:"author"`
}

type CurrencyAPIInput struct {
	Currency string `json:"currency"`
	Date     string `json:"date"`
}

type CurrencyAPIOutput struct {
	Service string             `json:"service"`
	Data    map[string]float64 `json:"data"`
}

func NewCurrencyAPIOutput() *CurrencyAPIOutput {
	return &CurrencyAPIOutput{
		Service: "currency",
		Data:    make(map[string]float64),
	}
}

func AppInfoHandler(appInfo InfoAPIOutput) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(appInfo); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
}

func addValuteToCurrencyAPIOutput(out *CurrencyAPIOutput, v cbr.CbrValute) error {
	val, err := cbr.ParseCbrValue(v.Value)
	if err != nil {
		log.Printf("failed to get value of currency: %v", err)
		return err
	}

	out.Data[v.CharCode] = val
	return nil
}

func fillCurrencyAPIOutput(out *CurrencyAPIOutput, curs cbr.CbrValCurs) error {
	for _, v := range curs.Valutes {
		err := addValuteToCurrencyAPIOutput(out, v)
		if err != nil {
			return err
		}
	}
	return nil
}

func AppCurrencyInfoHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var input CurrencyAPIInput
		input.Currency = req.URL.Query().Get("currency")
		input.Date = req.URL.Query().Get("date")

		log.Printf("currency %s, date %s\n", input.Currency, input.Date)

		values, err := cbr.GetCbrCurrencyValues(input.Date)
		if err != nil {
			http.Error(w, "failed to get currency values", http.StatusInternalServerError)
			return
		}

		output := NewCurrencyAPIOutput()
		if input.Currency == "" {
			fillCurrencyAPIOutput(output, values)
		} else {
			valute := cbr.FindCbrValute(values, input.Currency)
			if valute == nil {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}

			addValuteToCurrencyAPIOutput(output, *valute)
		}
		if err := json.NewEncoder(w).Encode(output); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
}
