package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// ExchangeRateAPI implements the API calls for https://exchangerate.host.
type ExchangeRateAPI struct {
	// url the API URL
	url string
}

type Currency struct {
	Motd struct {
		Msg string `json:"msg"`
		URL string `json:"url"`
	} `json:"motd"`
	Success bool `json:"success"`
	Query   struct {
		From   string `json:"from"`
		To     string `json:"to"`
		Amount int    `json:"amount"`
	} `json:"query"`
	Info struct {
		Rate float64 `json:"rate"`
	} `json:"info"`
	Historical bool    `json:"historical"`
	Date       string  `json:"date"`
	Result     float64 `json:"result"`
}

// NewExchangeRateAPI craetes a reference for exchangeRateApi.
func NewExchangeRateAPI(url string) *ExchangeRateAPI {
	return &ExchangeRateAPI{url: url}
}

// ConvertCurrency gets the current currency from a base currency.
func (e *ExchangeRateAPI) ConvertCurrency(from string, to string) (Currency, error) {
	resp, err := http.Get(fmt.Sprintf("%s/convert?from=%s&to=%s", e.url, url.QueryEscape(from), url.QueryEscape(to)))
	if err != nil {
		return Currency{}, fmt.Errorf("failed to call exchangerate api: %w", err)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Currency{}, fmt.Errorf("failed reading the request body: %w", err)
	}
	defer resp.Body.Close()

	currency := Currency{}

	err = json.Unmarshal(respBody, &currency)
	if err != nil {
		return Currency{}, fmt.Errorf("failed parsing the request body in JSON: %w", err)
	}

	return currency, nil
}
