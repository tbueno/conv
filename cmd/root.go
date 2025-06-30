package cmd

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"conv/internal/currency"
	"conv/internal/converter"
)

var rootCmd = &cobra.Command{
	Use:   "conv",
	Short: "A fast and simple CLI currency converter",
	Long: `A fast and simple CLI currency converter that supports a wide range of currencies.
Uses real-time exchange rates from Fawaz Ahmed's Currency API.

BACKWARD COMPATIBILITY:
  conv <amount> <from> <to>     # Direct conversion (legacy mode)
  conv --list                   # List currencies (legacy mode)

NEW SUBCOMMAND INTERFACE:
  conv convert <amount> <from> <to>   # Convert currencies
  conv list                           # List all currencies

Examples:
  conv 100 USD EUR              # Convert 100 USD to EUR (legacy)
  conv convert 100 USD EUR      # Convert 100 USD to EUR (new)
  conv --list                   # List currencies (legacy)  
  conv list                     # List currencies (new)`,
	Args: cobra.RangeArgs(0, 3),
	Run:  runRootCmd,
}

var listFlag bool

func init() {
	rootCmd.Flags().BoolVarP(&listFlag, "list", "l", false, "List all available currencies (legacy mode)")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func runRootCmd(cmd *cobra.Command, args []string) {
	// Handle --list flag (legacy mode)
	if listFlag {
		currency.ListCurrencies()
		return
	}

	// Handle direct conversion (legacy mode)
	if len(args) == 3 {
		// Parse and validate arguments
		input, err := ParseLegacyArgs(args)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
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
		return
	}

	// If no arguments provided, show help
	if len(args) == 0 {
		cmd.Help()
		return
	}

	// Invalid number of arguments
	fmt.Println("Error: invalid arguments")
	fmt.Println("Use 'conv --help' for usage information")
	os.Exit(1)
}

func ParseLegacyArgs(args []string) (currency.Input, error) {
	if len(args) != 3 {
		return currency.Input{}, fmt.Errorf("invalid number of arguments")
	}

	amount, err := strconv.ParseFloat(args[0], 32)
	if err != nil {
		return currency.Input{}, fmt.Errorf("invalid amount %s", args[0])
	}

	from := currency.Currency(strings.ToUpper(args[1]))
	to := currency.Currency(strings.ToUpper(args[2]))

	if !from.IsValid() {
		return currency.Input{}, fmt.Errorf("unsupported currency: %s", from)
	}
	if !to.IsValid() {
		return currency.Input{}, fmt.Errorf("unsupported currency: %s", to)
	}

	return currency.Input{
		Amount: float32(amount),
		From:   from,
		To:     to,
	}, nil
}