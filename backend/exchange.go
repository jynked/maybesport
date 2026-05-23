package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"sync"
	"time"
)

var (
	currentRate float64
	rateMu      sync.RWMutex

	lastSuccessfulRate float64
	lastRateMu         sync.RWMutex

	rateAvailable bool
)

const (
	defaultRate   = 11.5
	cacheDuration = 24 * time.Hour
)

func GetCurrentExchangeRate() float64 {
	rateMu.RLock()
	defer rateMu.RUnlock()

	if rateAvailable {
		return currentRate
	}
	lastRateMu.RLock()
	defer lastRateMu.RUnlock()
	if lastSuccessfulRate > 0 {
		return lastSuccessfulRate
	}
	return defaultRate
}

func refreshRateLoop() {
	updateExchangeRate()
	ticker := time.NewTicker(cacheDuration)
	go func() {
		for range ticker.C {
			updateExchangeRate()
		}
	}()
}

func updateExchangeRate() {
	if rate, err := fetchFromCbrXmlDaily(); err == nil && rate > 0 {
		setNewRate(rate, true)
		slog.Info("Exchange rate updated", "source", "cbr-xml-daily.ru", "rate", rate)
		return
	}
	if rate, err := fetchFromCoinGecko(); err == nil && rate > 0 {
		setNewRate(rate, true)
		slog.Info("Exchange rate updated", "source", "coingecko", "rate", rate)
		return
	}
	slog.Warn("All exchange rate APIs failed, using last successful rate", "last_rate", GetCurrentExchangeRate())
	setNewRate(GetCurrentExchangeRate(), false)
}

func setNewRate(rate float64, success bool) {
	if success {
		rateMu.Lock()
		currentRate = rate
		rateAvailable = true
		rateMu.Unlock()
		lastRateMu.Lock()
		lastSuccessfulRate = rate
		lastRateMu.Unlock()
		return
	}
	rateMu.Lock()
	rateAvailable = false
	rateMu.Unlock()
}

func fetchFromCbrXmlDaily() (float64, error) {
	url := "https://www.cbr-xml-daily.ru/daily_json.js"
	client := http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return 0, fmt.Errorf("cbr-xml-daily request failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("cbr-xml-daily bad status: %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("cbr-xml-daily read failed: %w", err)
	}
	var data struct {
		Valute map[string]struct {
			Value float64 `json:"Value"`
		} `json:"Valute"`
	}
	if err := json.Unmarshal(body, &data); err != nil {
		return 0, fmt.Errorf("cbr-xml-daily parse failed: %w", err)
	}
	cny, ok := data.Valute["CNY"]
	if !ok {
		return 0, fmt.Errorf("cbr-xml-daily: CNY not found")
	}
	if cny.Value <= 0 {
		return 0, fmt.Errorf("cbr-xml-daily: invalid CNY value")
	}
	return cny.Value, nil
}

func fetchFromCoinGecko() (float64, error) {
	url := "https://api.coingecko.com/api/v3/simple/price?ids=chinese-yuan&vs_currencies=rub"
	client := http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return 0, fmt.Errorf("coingecko request failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("coingecko bad status: %d", resp.StatusCode)
	}
	var data map[string]map[string]float64
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return 0, fmt.Errorf("coingecko parse failed: %w", err)
	}
	rate, ok := data["chinese-yuan"]["rub"]
	if !ok || rate == 0 {
		return 0, fmt.Errorf("coingecko: rate not found")
	}
	return rate, nil
}

func ExchangeRateHandler(w http.ResponseWriter, r *http.Request) {
	rate := GetCurrentExchangeRate()
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"rate":%.4f}`, rate)
}
