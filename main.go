package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Conversion struct {
	Date   string             `json:"date"`
	Values map[string]float32 `json:"-"`
}

// UnmarshalJSON implements custom JSON unmarshaling
func (c *Conversion) UnmarshalJSON(data []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	// Extract date
	if date, ok := raw["date"].(string); ok {
		c.Date = date
	}

	// Get the currency field (the only other field in the response)
	for key := range raw {
		if key == "date" {
			continue
		}
		if values, ok := raw[key].(map[string]interface{}); ok {
			c.Values = make(map[string]float32)
			for k, v := range values {
				if num, ok := v.(float64); ok {
					c.Values[k] = float32(num)
				}
			}
			return nil
		}
	}

	return fmt.Errorf("invalid response format: missing currency conversion map")
}

func convert(amount float32, from, to string, conversion *Conversion) (float32, error) {
	url := fmt.Sprintf(`https://cdn.jsdelivr.net/npm/@fawazahmed0/currency-api@latest/v1/currencies/%v.json`, from)

	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(conversion)
	if err != nil {
		return 0, err
	}

	if rate, exists := conversion.Values[to]; exists {
		return rate * amount, nil
	}
	return 0, fmt.Errorf("unsupported target currency: %s", to)
}

func listCurrencies() {
	var currencies map[string]string

	resp, err := http.Get("https://cdn.jsdelivr.net/npm/@fawazahmed0/currency-api@latest/v1/currencies.min.json")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&currencies)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("These are the supported currencies:")
	names := make([]string, 0, len(currencies))

	for _, name := range currencies {
		names = append(names, name)
	}

	sort.Strings(names)

	for _, name := range names {
		fmt.Printf("%s: %s\n", name, currencies[name])
	}
}

func main() {
	if len(os.Args) != 4 {
		fmt.Println("Usage: go run main.go <amount> <from> <to>")
		fmt.Println("Example: go run main.go 100 EUR USD")
		fmt.Println("Supported currencies: USD, EUR, BRL")
		os.Exit(1)
	}

	amount, err := strconv.ParseFloat(os.Args[1], 32)
	if err != nil {
		fmt.Println("Invalid amount. Please provide a valid number.")
		os.Exit(1)
	}

	from := strings.ToLower(os.Args[2])
	to := strings.ToLower(os.Args[3])

	conversion := &Conversion{}
	value, err := convert(float32(amount), from, to, conversion)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%v %s is %v %s\n", amount, strings.ToUpper(from), value, strings.ToUpper(to))
}
