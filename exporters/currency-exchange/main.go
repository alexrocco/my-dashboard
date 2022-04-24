package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var currencyExchangeRate = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "currency_exchange_rate",
		Help: "Current exchange rate.",
	},
	[]string{"from", "to"},
)

func init() {
	// Register the metric
	prometheus.MustRegister(currencyExchangeRate)
}

func main() {
	currecyFrom := os.Getenv("FROM")
	currencyTo := os.Getenv("TO")

	if len(currecyFrom) == 0 || len(currencyTo) == 0 {
		log.Fatal("environment variables 'FROM' and 'TO' are required")
	}

	exchangeRateAPI := NewExchangeRateAPI("https://api.exchangerate.host")

	ticker := time.NewTicker(30 * time.Second)
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				log.Println("Collecting currency exchange...")
				currency, err := exchangeRateAPI.ConvertCurrency(currecyFrom, currencyTo)
				if err != nil {
					log.Println(err)
					done <- true
				}

				currencyExchangeRate.With(prometheus.Labels{
					"from": currecyFrom,
					"to":   currencyTo,
				}).Set(currency.Result)
			}
		}
	}()

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":8080", nil))
}
