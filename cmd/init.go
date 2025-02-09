package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var force bool

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize teagit",
	Long: `Initialize teagit by creating a new profile.

This command sets up a new profile for teagit, including necessary configuration files.

Example:
	teagit init
	teagit init --force
`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := initializeTeagit(force); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("teagit initialized successfully!")
	},
}

func init() {
	initCmd.Flags().BoolVarP(&force, "force", "f", false, "Reinitialize teagit by overwriting existing configuration")
	rootCmd.AddCommand(initCmd)
}

// initializeTeagit handles the initialization logic
func initializeTeagit(force bool) error {
	configDir := filepath.Join(os.Getenv("HOME"), ".teagit")
	configFile := filepath.Join(configDir, "config.yaml")

	// Check if the configuration directory already exists
	if _, err := os.Stat(configDir); err == nil && !force {
		return fmt.Errorf("teagit is already initialized at %s. Use --force to overwrite", configDir)
	}

	// Remove existing config if force is set
	if force {
		fmt.Println("Force flag enabled. Removing existing configuration...")
		if err := os.RemoveAll(configDir); err != nil {
			return fmt.Errorf("failed to remove existing configuration: %v", err)
		}
	}

	// Create the configuration directory
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create configuration directory: %v", err)
	}
	fmt.Println("Configuration directory created successfully")

	// Set up Viper
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configDir)

	viper.SetDefault("system.init", true)
	viper.Set("system.init_time", time.Now().Format(time.RFC3339))

	// Write the configuration to a file
	fmt.Println("Writing configuration file:", configFile)
	if err := viper.WriteConfigAs(configFile); err != nil {
		return fmt.Errorf("failed to write configuration file: %v", err)
	}
	fmt.Println("Configuration file created successfully")

	return nil
}
