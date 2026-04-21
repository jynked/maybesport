package main

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"
)

type CBRResponse struct {
	Date         string `json:"Date"`
	PreviousDate string `json:"PreviousDate"`
	PreviousURL  string `json:"PreviousURL"`
	Timestamp    string `json:"Timestamp"`
	Valute       map[string]struct {
		ID       string  `json:"ID"`
		NumCode  string  `json:"NumCode"`
		CharCode string  `json:"CharCode"`
		Nominal  int     `json:"Nominal"`
		Name     string  `json:"Name"`
		Value    float64 `json:"Value"`
		Previous float64 `json:"Previous"`
	} `json:"Valute"`
}

var (
	exchangeRate   float64
	exchangeRateMu sync.RWMutex
)

const defaultExchangeRate = 11.5 // Запасной курс на случай недоступности ЦБ

func initExchangeRate() {
	rate, err := fetchExchangeRateFromCBR()
	if err != nil {
		exchangeRateMu.Lock()
		exchangeRate = defaultExchangeRate
		exchangeRateMu.Unlock()
	} else {
		exchangeRateMu.Lock()
		exchangeRate = rate
		exchangeRateMu.Unlock()
	}
}

func fetchExchangeRateFromCBR() (float64, error) {
	client := http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get("https://www.cbr-xml-daily.ru/daily_json.js")
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var data CBRResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return 0, err
	}

	cny, ok := data.Valute["CNY"]
	if !ok {
		return 0, nil
	}

	rate := cny.Value / float64(cny.Nominal)
	return rate, nil
}

func updateExchangeRate() {
	newRate, err := fetchExchangeRateFromCBR()
	if err != nil {
		return
	}
	exchangeRateMu.Lock()
	exchangeRate = newRate
	exchangeRateMu.Unlock()
}

func startExchangeRateUpdater() {
	updateExchangeRate()
	ticker := time.NewTicker(1 * time.Minute)
	go func() {
		for range ticker.C {
			updateExchangeRate()
		}
	}()
}
