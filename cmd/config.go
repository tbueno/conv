package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"conv/internal/config"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration settings",
	Long: `Manage configuration settings for the currency converter.

Available subcommands:
  set default-currency <CURRENCY>    Set the default target currency
  get default-currency               Show the current default currency
  show                               Show all configuration settings`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Error: config command requires a subcommand")
		fmt.Println("Use 'conv config --help' for usage information")
	},
}

var configSetCmd = &cobra.Command{
	Use:   "set <setting> <value>",
	Short: "Set a configuration value",
	Long: `Set a configuration value.

Available settings:
  default-currency <CURRENCY>    Set the default target currency

Examples:
  conv config set default-currency USD
  conv config set default-currency EUR`,
	Args: cobra.ExactArgs(2),
	Run:  runConfigSetCmd,
}

var configGetCmd = &cobra.Command{
	Use:   "get <setting>",
	Short: "Get a configuration value",
	Long: `Get a configuration value.

Available settings:
  default-currency    Show the current default currency

Examples:
  conv config get default-currency`,
	Args: cobra.ExactArgs(1),
	Run:  runConfigGetCmd,
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show all configuration settings",
	Long: `Show all configuration settings.

Examples:
  conv config show`,
	Args: cobra.NoArgs,
	Run:  runConfigShowCmd,
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configGetCmd)
	configCmd.AddCommand(configShowCmd)
}

func runConfigSetCmd(cmd *cobra.Command, args []string) {
	setting := strings.ToLower(args[0])
	value := args[1]

	switch setting {
	case "default-currency":
		err := config.SetDefaultCurrency(value)
		if err != nil {
			fmt.Printf("Error setting default currency: %v\n", err)
			return
		}
		fmt.Printf("Default currency set to: %s\n", strings.ToUpper(value))
	default:
		fmt.Printf("Error: unknown setting '%s'\n", setting)
		fmt.Println("Available settings: default-currency")
	}
}

func runConfigGetCmd(cmd *cobra.Command, args []string) {
	setting := strings.ToLower(args[0])

	switch setting {
	case "default-currency":
		currency, err := config.GetDefaultCurrency()
		if err != nil {
			fmt.Printf("Error getting default currency: %v\n", err)
			return
		}
		if currency == "" {
			fmt.Println("No default currency set")
		} else {
			fmt.Printf("Default currency: %s\n", currency)
		}
	default:
		fmt.Printf("Error: unknown setting '%s'\n", setting)
		fmt.Println("Available settings: default-currency")
	}
}

func runConfigShowCmd(cmd *cobra.Command, args []string) {
	cfg, err := config.GetConfig()
	if err != nil {
		fmt.Printf("Error loading configuration: %v\n", err)
		return
	}

	fmt.Println("Configuration:")
	if cfg.DefaultCurrency == "" {
		fmt.Println("  Default currency: (not set)")
	} else {
		fmt.Printf("  Default currency: %s\n", cfg.DefaultCurrency)
	}
}