package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
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

func listCurrencies() {
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

func parseArgs(args []string) (Input, error) {
	if len(args) != 4 || args[0] == "" || args[1] == "" || args[2] == "" || args[3] == "" {
		return Input{}, fmt.Errorf("invalid number of arguments")
	}

	amount, err := strconv.ParseFloat(args[1], 32)
	if err != nil {
		return Input{}, fmt.Errorf("invalid amount %s", args[1])
	}

	from := Currency(strings.ToUpper(args[2]))
	to := Currency(strings.ToUpper(args[3]))

	if !from.IsValid() {
		return Input{}, fmt.Errorf("unsupported currency: %s", from)
	}
	if !to.IsValid() {
		return Input{}, fmt.Errorf("unsupported currency: %s", to)
	}

	return Input{
		Amount: float32(amount),
		From:   from,
		To:     to,
	}, nil
}

func convert(input Input, conv Converter) (float32, error) {
	value, err := conv.Convert(input.Amount, strings.ToLower(input.From.String()), strings.ToLower(input.To.String()))
	if err != nil {
		return 0, err
	}
	return value, nil
}

func main() {
	// Check for --list flag
	if len(os.Args) == 2 && os.Args[1] == "--list" {
		listCurrencies()
		os.Exit(0)
	}

	input, err := parseArgs(os.Args)
	if err != nil {
		fmt.Println("Usage: go run main.go <amount> <from> <to>")
		fmt.Println("       go run main.go --list")
		fmt.Println("Example: go run main.go 100 EUR USD")
		fmt.Printf("Supported currencies: %s\n", strings.Join([]string{USD.String(), EUR.String(), BRL.String()}, ", "))
		os.Exit(1)
	}

	converter := &ApiCurrencyConverter{
		conversion: &FawazConversion{},
		apiUrl:     "https://cdn.jsdelivr.net/npm/@fawazahmed0/currency-api@latest/v1/currencies/%v.json",
	}

	value, err := convert(input, converter)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%v %s is %v %s\n", input.Amount, input.From, value, input.To)
}
