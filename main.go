package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
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

func (c Currency) String() string {
	return string(c)
}

func (c Currency) IsValid() bool {
	for _, supported := range supportedCurrencies {
		if c == supported {
			return true
		}
	}
	return false
}

func listCurrencies() {
	url := "https://cdn.jsdelivr.net/npm/@fawazahmed0/currency-api@latest/v1/currencies.min.json"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error fetching currencies: %v\n", err)
		fmt.Printf("Fallback - Supported currencies: %s\n", strings.Join([]string{USD.String(), EUR.String(), BRL.String()}, ", "))
		return
	}
	defer resp.Body.Close()

	var currencies map[string]string
	err = json.NewDecoder(resp.Body).Decode(&currencies)
	if err != nil {
		fmt.Printf("Error parsing currencies: %v\n", err)
		fmt.Printf("Fallback - Supported currencies: %s\n", strings.Join([]string{USD.String(), EUR.String(), BRL.String()}, ", "))
		return
	}

	fmt.Printf("Available currencies (%d total):\n", len(currencies))
	for code, name := range currencies {
		fmt.Printf("  %s - %s\n", strings.ToUpper(code), name)
	}
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
