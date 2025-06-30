package cmd

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"conv/internal/currency"
	"conv/internal/converter"
)

var convertCmd = &cobra.Command{
	Use:   "convert <amount> <from> <to>",
	Short: "Convert currency amounts between different currencies",
	Long: `Convert currency amounts between different currencies using real-time exchange rates.

Examples:
  conv convert 100 USD EUR    # Convert 100 USD to EUR
  conv convert 50 GBP JPY     # Convert 50 GBP to JPY
  conv convert 1000 BTC USD   # Convert 1000 BTC to USD`,
	Args: cobra.MatchAll(cobra.ExactArgs(3), validateConvertArgs),
	Run:  runConvertCmd,
}

func init() {
	rootCmd.AddCommand(convertCmd)
}

func validateConvertArgs(cmd *cobra.Command, args []string) error {
	if len(args) != 3 {
		return fmt.Errorf("requires exactly 3 arguments: <amount> <from> <to>")
	}

	// Validate amount
	if _, err := strconv.ParseFloat(args[0], 32); err != nil {
		return fmt.Errorf("invalid amount '%s': must be a valid number", args[0])
	}

	// Validate currencies
	from := currency.Currency(strings.ToUpper(args[1]))
	to := currency.Currency(strings.ToUpper(args[2]))

	if !from.IsValid() {
		return fmt.Errorf("unsupported source currency: %s", from)
	}
	if !to.IsValid() {
		return fmt.Errorf("unsupported target currency: %s", to)
	}

	return nil
}

func runConvertCmd(cmd *cobra.Command, args []string) {
	// Arguments are already validated by Cobra, so we can safely parse them
	input := parseConvertArgs(args)

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

func parseConvertArgs(args []string) currency.Input {
	// Args are pre-validated, so we can safely parse without error checking
	amount, _ := strconv.ParseFloat(args[0], 32)
	from := currency.Currency(strings.ToUpper(args[1]))
	to := currency.Currency(strings.ToUpper(args[2]))

	return currency.Input{
		Amount: float32(amount),
		From:   from,
		To:     to,
	}
}