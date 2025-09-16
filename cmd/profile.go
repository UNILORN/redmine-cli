package cmd

import (
	"fmt"
	"redmine-cli/config"

	"github.com/spf13/cobra"
)

var profileCmd = &cobra.Command{
	Use:   "profile",
	Short: "Manage profiles",
	Long:  `Manage multiple Redmine server profiles`,
}

var profileAddCmd = &cobra.Command{
	Use:   "add [name] [url] [token]",
	Short: "Add a new profile",
	Long:  `Add a new profile with name, Redmine URL, and API token`,
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		url := args[1]
		token := args[2]

		cfg, err := config.Load()
		if err != nil {
			fmt.Printf("Error loading config: %v\n", err)
			return
		}

		if err := cfg.AddProfile(name, url, token); err != nil {
			fmt.Printf("Error adding profile: %v\n", err)
			return
		}

		if err := cfg.Save(); err != nil {
			fmt.Printf("Error saving config: %v\n", err)
			return
		}

		fmt.Printf("Profile '%s' has been added successfully\n", name)
		if cfg.DefaultProfile == name {
			fmt.Printf("Set as default profile\n")
		}
	},
}

var profileListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all profiles",
	Long:  `List all configured profiles`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load()
		if err != nil {
			fmt.Printf("Error loading config: %v\n", err)
			return
		}

		if len(cfg.Profiles) == 0 {
			fmt.Println("No profiles configured.")
			return
		}

		fmt.Printf("Configured profiles:\n\n")
		for name, profile := range cfg.Profiles {
			marker := "  "
			if name == cfg.DefaultProfile {
				marker = "* "
			}
			fmt.Printf("%s%s\n", marker, name)
			fmt.Printf("  URL: %s\n", profile.RedmineURL)
			fmt.Printf("  API Key: %s\n", maskAPIKey(profile.APIKey))
			fmt.Println()
		}

		if cfg.DefaultProfile != "" {
			fmt.Printf("* Default profile\n")
		}
	},
}

var profileUseCmd = &cobra.Command{
	Use:   "use [name]",
	Short: "Set default profile",
	Long:  `Set the default profile to use`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		cfg, err := config.Load()
		if err != nil {
			fmt.Printf("Error loading config: %v\n", err)
			return
		}

		if err := cfg.SetDefaultProfile(name); err != nil {
			fmt.Printf("Error setting default profile: %v\n", err)
			return
		}

		if err := cfg.Save(); err != nil {
			fmt.Printf("Error saving config: %v\n", err)
			return
		}

		fmt.Printf("Default profile set to '%s'\n", name)
	},
}

var profileRemoveCmd = &cobra.Command{
	Use:   "remove [name]",
	Short: "Remove a profile",
	Long:  `Remove a profile from configuration`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		cfg, err := config.Load()
		if err != nil {
			fmt.Printf("Error loading config: %v\n", err)
			return
		}

		if err := cfg.RemoveProfile(name); err != nil {
			fmt.Printf("Error removing profile: %v\n", err)
			return
		}

		if err := cfg.Save(); err != nil {
			fmt.Printf("Error saving config: %v\n", err)
			return
		}

		fmt.Printf("Profile '%s' has been removed\n", name)
	},
}

var profileShowCmd = &cobra.Command{
	Use:   "show [name]",
	Short: "Show profile details",
	Long:  `Show detailed information about a specific profile`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load()
		if err != nil {
			fmt.Printf("Error loading config: %v\n", err)
			return
		}

		var profileName string
		if len(args) == 1 {
			profileName = args[0]
		} else {
			profileName = cfg.DefaultProfile
			if profileName == "" {
				fmt.Println("No default profile set and no profile specified.")
				return
			}
		}

		profile, exists := cfg.Profiles[profileName]
		if !exists {
			fmt.Printf("Profile '%s' not found\n", profileName)
			return
		}

		fmt.Printf("Profile: %s\n", profileName)
		if profileName == cfg.DefaultProfile {
			fmt.Println("Status: Default profile")
		}
		fmt.Printf("Redmine URL: %s\n", profile.RedmineURL)
		fmt.Printf("API Key: %s\n", maskAPIKey(profile.APIKey))
	},
}

func init() {
	rootCmd.AddCommand(profileCmd)
	profileCmd.AddCommand(profileAddCmd)
	profileCmd.AddCommand(profileListCmd)
	profileCmd.AddCommand(profileUseCmd)
	profileCmd.AddCommand(profileRemoveCmd)
	profileCmd.AddCommand(profileShowCmd)
}