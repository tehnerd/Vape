package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tehnerd/vape/internal/config"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration",
	Long:  "Manage VAPE configuration settings",
}

var configInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize configuration file",
	Long:  "Create a new configuration file with default settings",
	RunE: func(cmd *cobra.Command, args []string) error {
		if config.ConfigFileExists() {
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("Config file already exists. Overwrite? [y/N]: ")
			response, _ := reader.ReadString('\n')
			response = strings.TrimSpace(strings.ToLower(response))
			if response != "y" && response != "yes" {
				fmt.Println("Aborted.")
				return nil
			}
		}

		reader := bufio.NewReader(os.Stdin)

		fmt.Print("Enter your RIPE Atlas API key: ")
		apiKey, _ := reader.ReadString('\n')
		apiKey = strings.TrimSpace(apiKey)

		fmt.Print("Default output format (table/json) [table]: ")
		outputFormat, _ := reader.ReadString('\n')
		outputFormat = strings.TrimSpace(outputFormat)
		if outputFormat == "" {
			outputFormat = "table"
		}

		fmt.Print("Default address family (4/6) [4]: ")
		afStr, _ := reader.ReadString('\n')
		afStr = strings.TrimSpace(afStr)
		af := 4
		if afStr == "6" {
			af = 6
		}

		fmt.Print("Default number of probes [10]: ")
		probesStr, _ := reader.ReadString('\n')
		probesStr = strings.TrimSpace(probesStr)
		probes := 10
		if probesStr != "" {
			fmt.Sscanf(probesStr, "%d", &probes)
		}

		if err := config.WriteConfig(apiKey, outputFormat, af, probes); err != nil {
			return fmt.Errorf("failed to write config: %w", err)
		}

		fmt.Printf("Config written to %s\n", config.GetConfigPath())
		return nil
	},
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	Long:  "Display current configuration settings",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Config file: %s\n", config.GetConfigPath())
		fmt.Printf("Config exists: %t\n\n", config.ConfigFileExists())
		fmt.Printf("API Key: %s\n", maskAPIKey(config.GetAPIKey()))
		fmt.Printf("Output Format: %s\n", config.GetOutputFormat())
		fmt.Printf("Default AF: %d\n", config.GetDefaultAF())
		fmt.Printf("Default Probes: %d\n", config.GetDefaultProbes())
	},
}

var configSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set a configuration value",
	Long: `Set a configuration value. Available keys:
  - api_key: Your RIPE Atlas API key
  - output_format: Default output format (table/json)
  - default_af: Default address family (4/6)
  - default_probes: Default number of probes`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		key := args[0]
		value := args[1]

		validKeys := map[string]bool{
			"api_key":        true,
			"output_format":  true,
			"default_af":     true,
			"default_probes": true,
		}

		if !validKeys[key] {
			return fmt.Errorf("invalid key: %s", key)
		}

		viper.Set(key, value)
		if err := viper.WriteConfigAs(config.GetConfigPath()); err != nil {
			return fmt.Errorf("failed to write config: %w", err)
		}

		fmt.Printf("Set %s = %s\n", key, value)
		return nil
	},
}

func maskAPIKey(key string) string {
	if key == "" {
		return "(not set)"
	}
	if len(key) <= 8 {
		return "****"
	}
	return key[:4] + "..." + key[len(key)-4:]
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configInitCmd)
	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configSetCmd)
}
