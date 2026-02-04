package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tehnerd/vape/cmd/credits"
	"github.com/tehnerd/vape/cmd/measurements"
	"github.com/tehnerd/vape/cmd/probes"
	"github.com/tehnerd/vape/cmd/quick"
	"github.com/tehnerd/vape/internal/config"
)

var (
	cfgFile      string
	apiKey       string
	outputFormat string
)

var rootCmd = &cobra.Command{
	Use:   "vape",
	Short: "VAPE - RIPE Atlas CLI Tool",
	Long: `VAPE is a command-line interface for interacting with the RIPE Atlas API.

It allows you to create and manage network measurements, query probes,
check credits, and run quick one-off tests.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.vape.yaml)")
	rootCmd.PersistentFlags().StringVar(&apiKey, "key", "", "RIPE Atlas API key")
	rootCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "table", "output format (table, json)")

	viper.BindPFlag("api_key", rootCmd.PersistentFlags().Lookup("key"))
	viper.BindPFlag("output_format", rootCmd.PersistentFlags().Lookup("output"))
}

func initConfig() {
	config.InitConfig(cfgFile)
}

func init() {
	rootCmd.AddCommand(measurements.MeasurementsCmd)
	rootCmd.AddCommand(probes.ProbesCmd)
	rootCmd.AddCommand(credits.CreditsCmd)
	rootCmd.AddCommand(quick.QuickCmd)
}
