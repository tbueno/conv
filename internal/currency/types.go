package currency

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type Currency string

const (
	USD Currency = "USD"
	EUR Currency = "EUR"
	BRL Currency = "BRL"
)

var supportedCurrencies = []Currency{USD, EUR, BRL}
var cachedCurrencies map[string]string

func (c Currency) String() string {
	return string(c)
}

func (c Currency) IsValid() bool {
	// Load cached currencies if not already loaded
	if cachedCurrencies == nil {
		loadCachedCurrencies()
	}
	
	// Check cached currencies first
	if cachedCurrencies != nil {
		_, exists := cachedCurrencies[strings.ToLower(string(c))]
		return exists
	}
	
	// Fallback to hardcoded currencies
	for _, supported := range supportedCurrencies {
		if c == supported {
			return true
		}
	}
	return false
}

func getCacheFilePath() string {
	return filepath.Join("conf", "currencies.json")
}

func loadCachedCurrencies() {
	cacheFile := getCacheFilePath()
	data, err := os.ReadFile(cacheFile)
	if err != nil {
		// Cache file doesn't exist, try to download and cache
		downloadAndCacheCurrencies()
		return
	}
	
	err = json.Unmarshal(data, &cachedCurrencies)
	if err != nil {
		log.Printf("Error parsing cached currencies: %v", err)
		cachedCurrencies = nil
	}
}

func downloadAndCacheCurrencies() {
	url := "https://cdn.jsdelivr.net/npm/@fawazahmed0/currency-api@latest/v1/currencies.min.json"
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Error fetching currencies for cache: %v", err)
		return
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading currency data: %v", err)
		return
	}

	err = json.Unmarshal(data, &cachedCurrencies)
	if err != nil {
		log.Printf("Error parsing currency data: %v", err)
		return
	}

	// Create conf directory if it doesn't exist
	confDir := filepath.Dir(getCacheFilePath())
	err = os.MkdirAll(confDir, 0755)
	if err != nil {
		log.Printf("Error creating conf directory: %v", err)
		return
	}

	// Save to cache file
	err = os.WriteFile(getCacheFilePath(), data, 0644)
	if err != nil {
		log.Printf("Error saving currencies to cache: %v", err)
	}
}

func ListCurrencies() {
	// Load cached currencies if not already loaded
	if cachedCurrencies == nil {
		loadCachedCurrencies()
	}
	
	if cachedCurrencies != nil {
		fmt.Printf("Available currencies (%d total):\n", len(cachedCurrencies))
		for code, name := range cachedCurrencies {
			fmt.Printf("  %s - %s\n", strings.ToUpper(code), name)
		}
		return
	}
	
	// Fallback if caching failed
	fmt.Printf("Fallback - Supported currencies: %s\n", strings.Join([]string{USD.String(), EUR.String(), BRL.String()}, ", "))
}

type Input struct {
	Amount float32
	From   Currency
	To     Currency
}