package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

const (
	DefaultOutputFormat = "table"
	DefaultAF           = 4
	DefaultProbes       = 10
)

func InitConfig(cfgFile string) {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error finding home directory:", err)
			os.Exit(1)
		}

		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".vape")
	}

	// Environment variables
	viper.SetEnvPrefix("")
	viper.BindEnv("api_key", "RIPE_ATLAS_KEY")
	viper.AutomaticEnv()

	// Set defaults
	viper.SetDefault("output_format", DefaultOutputFormat)
	viper.SetDefault("default_af", DefaultAF)
	viper.SetDefault("default_probes", DefaultProbes)

	// Read config file if it exists
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			fmt.Fprintln(os.Stderr, "Error reading config file:", err)
		}
	}
}

func GetAPIKey() string {
	return viper.GetString("api_key")
}

func GetOutputFormat() string {
	return viper.GetString("output_format")
}

func GetDefaultAF() int {
	return viper.GetInt("default_af")
}

func GetDefaultProbes() int {
	return viper.GetInt("default_probes")
}

func GetConfigPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".vape.yaml")
}

func ConfigFileExists() bool {
	path := GetConfigPath()
	_, err := os.Stat(path)
	return err == nil
}

func WriteConfig(apiKey, outputFormat string, defaultAF, defaultProbes int) error {
	viper.Set("api_key", apiKey)
	viper.Set("output_format", outputFormat)
	viper.Set("default_af", defaultAF)
	viper.Set("default_probes", defaultProbes)

	path := GetConfigPath()
	return viper.WriteConfigAs(path)
}
