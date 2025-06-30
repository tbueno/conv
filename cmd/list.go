package cmd

import (
	"github.com/spf13/cobra"
	"conv/internal/currency"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available currencies",
	Long: `List all available currencies supported by the converter.

This includes fiat currencies, cryptocurrencies, and precious metals.
The list is cached locally for faster performance.

Examples:
  conv list               # List all currencies
  conv list | grep USD    # Search for USD-related currencies`,
	Args: cobra.NoArgs,
	Run:  runListCmd,
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func runListCmd(cmd *cobra.Command, args []string) {
	currency.ListCurrencies()
}