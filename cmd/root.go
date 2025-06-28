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
	Use:   "conv <amount> <from> <to>",
	Short: "Convert currency amounts between different currencies",
	Long: `A fast and simple CLI currency converter that supports a wide range of currencies.
Uses real-time exchange rates from Fawaz Ahmed's Currency API.

Examples:
  conv 100 USD EUR    # Convert 100 USD to EUR
  conv 50 GBP JPY     # Convert 50 GBP to JPY
  conv --list         # List all available currencies`,
	Args: cobra.RangeArgs(0, 3),
	Run:  runConvert,
}

var listFlag bool

func init() {
	rootCmd.Flags().BoolVarP(&listFlag, "list", "l", false, "List all available currencies")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func runConvert(cmd *cobra.Command, args []string) {
	// Handle --list flag
	if listFlag {
		currency.ListCurrencies()
		return
	}

	// Check if we have the required arguments for conversion
	if len(args) != 3 {
		fmt.Println("Error: conversion requires exactly 3 arguments: <amount> <from> <to>")
		fmt.Println("Use 'conv --help' for more information")
		os.Exit(1)
	}

	// Parse and validate arguments
	input, err := parseArgsFromCobra(args)
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
}

func parseArgsFromCobra(args []string) (currency.Input, error) {
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