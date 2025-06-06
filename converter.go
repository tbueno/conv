package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Converter interface {
	Convert(amount float32, from, to string) (float32, error)
}

type FawazConversion struct {
	Date   string             `json:"date"`
	Values map[string]float32 `json:"-"`
}

type ApiCurrencyConverter struct {
	conversion *FawazConversion
	apiUrl     string
}

func (c *ApiCurrencyConverter) Convert(amount float32, from, to string) (float32, error) {
	url := fmt.Sprintf(c.apiUrl, from)
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(c.conversion)
	if err != nil {
		return 0, err
	}

	if rate, exists := c.conversion.Values[to]; exists {
		return rate * amount, nil
	}
	return 0, fmt.Errorf("unsupported target currency: %s", to)
}

// UnmarshalJSON implements custom JSON unmarshaling
func (c *FawazConversion) UnmarshalJSON(data []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

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
