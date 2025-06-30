package currency

import (
	"testing"
)

func TestCurrencyValidation(t *testing.T) {
	tests := []struct {
		name     string
		currency Currency
		want     bool
	}{
		{
			name:     "hardcoded USD should be valid",
			currency: USD,
			want:     true,
		},
		{
			name:     "hardcoded EUR should be valid",
			currency: EUR,
			want:     true,
		},
		{
			name:     "JPY should be valid (from API cache)",
			currency: "JPY",
			want:     true,
		},
		{
			name:     "invalid currency should be false",
			currency: "INVALID",
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.currency.IsValid()
			if got != tt.want {
				t.Errorf("Currency.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCurrencyString(t *testing.T) {
	tests := []struct {
		name     string
		currency Currency
		want     string
	}{
		{
			name:     "USD to string",
			currency: USD,
			want:     "USD",
		},
		{
			name:     "EUR to string",
			currency: EUR,
			want:     "EUR",
		},
		{
			name:     "Custom currency to string",
			currency: "BTC",
			want:     "BTC",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.currency.String()
			if got != tt.want {
				t.Errorf("Currency.String() = %v, want %v", got, tt.want)
			}
		})
	}
}