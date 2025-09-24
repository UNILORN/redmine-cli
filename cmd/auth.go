package cmd

import (
	"fmt"
	"strings"

	"github.com/UNILORN/redmine-cli/config"

	"github.com/spf13/cobra"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authentication management",
	Long:  `Manage authentication tokens and credentials for Redmine`,
}

var tokenCmd = &cobra.Command{
	Use:   "token",
	Short: "Token management",
	Long:  `Manage API tokens for Redmine authentication`,
}

var tokenAddCmd = &cobra.Command{
	Use:   "add [token]",
	Short: "Add API token to current profile",
	Long:  `Add an API token to the current default profile`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		token := args[0]

		cfg, err := config.Load()
		if err != nil {
			fmt.Printf("Error loading config: %v\n", err)
			return
		}

		if cfg.DefaultProfile == "" {
			fmt.Println("No default profile configured. Please add a profile first using 'redmine profile add'")
			return
		}

		profile, exists := cfg.Profiles[cfg.DefaultProfile]
		if !exists {
			fmt.Printf("Default profile '%s' not found\n", cfg.DefaultProfile)
			return
		}

		profile.APIKey = token
		cfg.Profiles[cfg.DefaultProfile] = profile

		if err := cfg.Save(); err != nil {
			fmt.Printf("Error saving config: %v\n", err)
			return
		}

		fmt.Printf("API token has been saved to profile '%s' successfully\n", cfg.DefaultProfile)
	},
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configuration management",
	Long:  `Manage Redmine CLI configuration (deprecated - use 'profile' command instead)`,
}

var setURLCmd = &cobra.Command{
	Use:   "set-url [url]",
	Short: "Set Redmine URL for current profile",
	Long:  `Set the Redmine server URL for the current default profile (deprecated - use 'profile' command instead)`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		url := args[0]

		cfg, err := config.Load()
		if err != nil {
			fmt.Printf("Error loading config: %v\n", err)
			return
		}

		if cfg.DefaultProfile == "" {
			fmt.Println("No default profile configured. Please add a profile first using 'redmine profile add'")
			return
		}

		profile, exists := cfg.Profiles[cfg.DefaultProfile]
		if !exists {
			fmt.Printf("Default profile '%s' not found\n", cfg.DefaultProfile)
			return
		}

		profile.RedmineURL = url
		cfg.Profiles[cfg.DefaultProfile] = profile

		if err := cfg.Save(); err != nil {
			fmt.Printf("Error saving config: %v\n", err)
			return
		}

		fmt.Printf("Redmine URL has been saved to profile '%s' successfully\n", cfg.DefaultProfile)
	},
}

var showConfigCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	Long:  `Show current configuration (deprecated - use 'profile show' instead)`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load()
		if err != nil {
			fmt.Printf("Error loading config: %v\n", err)
			return
		}

		if cfg.DefaultProfile == "" || len(cfg.Profiles) == 0 {
			fmt.Println("No profiles configured. Use 'redmine profile add' to create a profile.")
			return
		}

		profile, exists := cfg.Profiles[cfg.DefaultProfile]
		if !exists {
			fmt.Printf("Default profile '%s' not found\n", cfg.DefaultProfile)
			return
		}

		fmt.Printf("Current profile: %s\n", cfg.DefaultProfile)
		fmt.Printf("Redmine URL: %s\n", profile.RedmineURL)

		if profile.APIKey != "" {
			fmt.Printf("API Key: %s\n", maskAPIKey(profile.APIKey))
		} else {
			fmt.Println("API Key: Not configured")
		}
	},
}

func maskAPIKey(key string) string {
	if len(key) <= 8 {
		return strings.Repeat("*", len(key))
	}
	return key[:4] + strings.Repeat("*", len(key)-8) + key[len(key)-4:]
}

func init() {
	rootCmd.AddCommand(authCmd)
	rootCmd.AddCommand(configCmd)

	authCmd.AddCommand(tokenCmd)
	tokenCmd.AddCommand(tokenAddCmd)

	configCmd.AddCommand(setURLCmd)
	configCmd.AddCommand(showConfigCmd)
}
