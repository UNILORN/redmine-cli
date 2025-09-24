package cmd

import (
	"fmt"
	"redmine-cli/config"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var urlIssueCmd = &cobra.Command{
	Use:   "url [issue_id]",
	Short: "Get the URL for an issue",
	Long:  `Get the URL for a specific issue in Redmine.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load()
		if err != nil {
			fmt.Printf("Error loading config: %v\n", err)
			return
		}

		profile, err := cfg.GetCurrentProfile()
		if err != nil {
			fmt.Printf("Error getting current profile: %v\n", err)
			fmt.Println("Please add a profile using 'redmine profile add'")
			return
		}

		if profile.RedmineURL == "" {
			fmt.Printf("Redmine URL not configured for profile '%s'. Please run 'redmine profile add'\n", profile.Name)
			return
		}

		// Parse issue ID
		issueID, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Printf("Invalid issue ID: %s\n", args[0])
			return
		}

		// Remove trailing slash from RedmineURL if present
		baseURL := strings.TrimSuffix(profile.RedmineURL, "/")

		// Print the issue URL
		fmt.Printf("%s/issues/%d\n", baseURL, issueID)
	},
}
