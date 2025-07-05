package cmd

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"conv/internal/config"
	"conv/internal/currency"
	"conv/internal/converter"
)

var convertCmd = &cobra.Command{
	Use:   "convert <amount> <from> [to]",
	Short: "Convert currency amounts between different currencies",
	Long: `Convert currency amounts between different currencies using real-time exchange rates.

If no target currency is specified, the default currency will be used.
Set a default currency with: conv config set default-currency <CURRENCY>

Examples:
  conv convert 100 USD EUR    # Convert 100 USD to EUR
  conv convert 100 USD        # Convert 100 USD to default currency
  conv convert 50 GBP JPY     # Convert 50 GBP to JPY
  conv convert 1000 BTC USD   # Convert 1000 BTC to USD`,
	Args: cobra.MatchAll(cobra.RangeArgs(2, 3), validateConvertArgs),
	Run:  runConvertCmd,
}

func init() {
	rootCmd.AddCommand(convertCmd)
}

func validateConvertArgs(cmd *cobra.Command, args []string) error {
	if len(args) < 2 || len(args) > 3 {
		return fmt.Errorf("requires 2 or 3 arguments: <amount> <from> [to]")
	}

	// Validate amount
	if _, err := strconv.ParseFloat(args[0], 32); err != nil {
		return fmt.Errorf("invalid amount '%s': must be a valid number", args[0])
	}

	// Validate source currency
	from := currency.Currency(strings.ToUpper(args[1]))
	if !from.IsValid() {
		return fmt.Errorf("unsupported source currency: %s", from)
	}

	// Validate target currency if provided
	if len(args) == 3 {
		to := currency.Currency(strings.ToUpper(args[2]))
		if !to.IsValid() {
			return fmt.Errorf("unsupported target currency: %s", to)
		}
	} else {
		// If no target currency provided, check if default currency is set
		defaultCurrency, err := config.GetDefaultCurrency()
		if err != nil {
			return fmt.Errorf("failed to get default currency: %v", err)
		}
		if defaultCurrency == "" {
			return fmt.Errorf("no target currency specified and no default currency set. Use 'conv config set default-currency <CURRENCY>' to set a default")
		}
	}

	return nil
}

func runConvertCmd(cmd *cobra.Command, args []string) {
	// Arguments are already validated by Cobra, so we can safely parse them
	input, err := parseConvertArgs(args)
	if err != nil {
		log.Fatal(err)
	}

	// Perform conversion
	conv := &converter.ApiCurrencyConverter{
		Conversion: &converter.FawazConversion{},
		ApiUrl:     "https://cdn.jsdelivr.net/npm/@fawazahmed0/currency-api@latest/v1/currencies/%v.json",
	}

	value, err := converter.Convert(input, conv)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%v %s is %v %s\n", input.Amount, input.From, value, input.To)
}

func parseConvertArgs(args []string) (currency.Input, error) {
	// Args are pre-validated, so we can safely parse without error checking
	amount, _ := strconv.ParseFloat(args[0], 32)
	from := currency.Currency(strings.ToUpper(args[1]))
	
	var to currency.Currency
	if len(args) == 3 {
		// Target currency explicitly provided
		to = currency.Currency(strings.ToUpper(args[2]))
	} else {
		// Use default currency
		defaultCurrency, err := config.GetDefaultCurrency()
		if err != nil {
			return currency.Input{}, fmt.Errorf("failed to get default currency: %v", err)
		}
		to = defaultCurrency
	}

	return currency.Input{
		Amount: float32(amount),
		From:   from,
		To:     to,
	}, nil
}