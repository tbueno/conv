package converter

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"conv/internal/currency"
)

type Converter interface {
	Convert(amount float32, from, to string) (float32, error)
}

type FawazConversion struct {
	Date   string             `json:"date"`
	Values map[string]float32 `json:"-"`
}

type ApiCurrencyConverter struct {
	Conversion *FawazConversion
	ApiUrl     string
}

func (c *ApiCurrencyConverter) Convert(amount float32, from, to string) (float32, error) {
	url := fmt.Sprintf(c.ApiUrl, from)
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(c.Conversion)
	if err != nil {
		return 0, err
	}

	if rate, exists := c.Conversion.Values[to]; exists {
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

func Convert(input currency.Input, conv Converter) (float32, error) {
	value, err := conv.Convert(input.Amount, strings.ToLower(input.From.String()), strings.ToLower(input.To.String()))
	if err != nil {
		return 0, err
	}
	return value, nil
}