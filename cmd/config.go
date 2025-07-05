package cmd

import (
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
		cmd.Println("Error: config command requires a subcommand")
		cmd.Println("Use 'conv config --help' for usage information")
	},
}

var configSetCmd = &cobra.Command{
	Use:   "set <setting> <value>",
	Short: "Set a configuration value",
	Long: `Set a configuration value.

Available settings:
  default-currency <CURRENCY>    Set the default target currency
  default-currency clear         Clear the default target currency

Examples:
  conv config set default-currency USD
  conv config set default-currency EUR
  conv config set default-currency clear`,
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
		if value == "" || strings.ToLower(value) == "none" || strings.ToLower(value) == "clear" {
			err := config.ClearDefaultCurrency()
			if err != nil {
				cmd.Printf("Error clearing default currency: %v\n", err)
				return
			}
			cmd.Println("Default currency cleared")
		} else {
			err := config.SetDefaultCurrency(value)
			if err != nil {
				cmd.Printf("Error setting default currency: %v\n", err)
				return
			}
			cmd.Printf("Default currency set to: %s\n", strings.ToUpper(value))
		}
	default:
		cmd.Printf("Error: unknown setting '%s'\n", setting)
		cmd.Println("Available settings: default-currency")
	}
}

func runConfigGetCmd(cmd *cobra.Command, args []string) {
	setting := strings.ToLower(args[0])

	switch setting {
	case "default-currency":
		currency, err := config.GetDefaultCurrency()
		if err != nil {
			cmd.Printf("Error getting default currency: %v\n", err)
			return
		}
		if currency == "" {
			cmd.Println("No default currency set")
		} else {
			cmd.Printf("Default currency: %s\n", currency)
		}
	default:
		cmd.Printf("Error: unknown setting '%s'\n", setting)
		cmd.Println("Available settings: default-currency")
	}
}

func runConfigShowCmd(cmd *cobra.Command, args []string) {
	cfg, err := config.GetConfig()
	if err != nil {
		cmd.Printf("Error loading configuration: %v\n", err)
		return
	}

	cmd.Println("Configuration:")
	if cfg.DefaultCurrency == "" {
		cmd.Println("  Default currency: (not set)")
	} else {
		cmd.Printf("  Default currency: %s\n", cfg.DefaultCurrency)
	}
}