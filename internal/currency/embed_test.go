package currency

import (
	"encoding/json"
	"testing"
)

func TestEmbeddedCurrencies(t *testing.T) {
	// Test that the embedded file exists and is valid JSON
	data, err := embeddedFiles.ReadFile("conf/currencies.json")
	if err != nil {
		t.Fatalf("Failed to read embedded currencies.json: %v", err)
	}

	// Test that it's valid JSON
	var currencies map[string]string
	err = json.Unmarshal(data, &currencies)
	if err != nil {
		t.Fatalf("Failed to parse embedded currencies.json: %v", err)
	}

	// Test that it contains expected currencies
	expectedCurrencies := []string{"usd", "eur", "gbp", "jpy", "btc", "eth"}
	for _, curr := range expectedCurrencies {
		if _, exists := currencies[curr]; !exists {
			t.Errorf("Expected currency %s not found in embedded file", curr)
		}
	}

	// Test that it's not empty
	if len(currencies) == 0 {
		t.Error("Embedded currencies file is empty")
	}

	t.Logf("Successfully loaded %d currencies from embedded file", len(currencies))
}

func TestLoadingPriority(t *testing.T) {
	// Reset cached currencies to test loading priority
	cachedCurrencies = nil

	// This should load from embedded file first
	loadCachedCurrencies()

	if cachedCurrencies == nil {
		t.Fatal("Failed to load currencies from any source")
	}

	// Verify we have a reasonable number of currencies
	if len(cachedCurrencies) < 100 {
		t.Errorf("Expected at least 100 currencies, got %d", len(cachedCurrencies))
	}

	t.Logf("Loaded %d currencies", len(cachedCurrencies))
}