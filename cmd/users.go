package cmd

import (
	"fmt"
	"redmine-cli/client"
	"redmine-cli/config"
	"strings"

	"github.com/spf13/cobra"
)

var usersCmd = &cobra.Command{
	Use:   "users",
	Short: "Manage Redmine users",
	Long:  `List and view Redmine users`,
}

var listUsersCmd = &cobra.Command{
	Use:   "list",
	Short: "List users",
	Long:  `List all users from Redmine`,
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

		if profile.APIKey == "" {
			fmt.Printf("API key not configured for profile '%s'. Please run 'redmine auth token add <token>' or 'redmine profile add'\n", profile.Name)
			return
		}

		if profile.RedmineURL == "" {
			fmt.Printf("Redmine URL not configured for profile '%s'. Please run 'redmine profile add'\n", profile.Name)
			return
		}

		c := client.NewClient(profile.RedmineURL, profile.APIKey)

		response, err := c.GetUsers()
		if err != nil {
			fmt.Printf("Error getting users: %v\n", err)
			return
		}

		if len(response.Users) == 0 {
			fmt.Println("No users found.")
			return
		}

		fmt.Printf("Users (Total: %d)\n", len(response.Users))
		fmt.Println(strings.Repeat("-", 80))

		for _, user := range response.Users {
			name := user.Name
			if name == "" {
				name = "(No name)"
			}
			login := user.Login
			if login == "" {
				login = "(No login)"
			}
			email := user.Email
			if email == "" {
				email = "(No email)"
			}
			fmt.Printf("ID: %d | Name: %s | Login: %s | Email: %s\n", user.ID, name, login, email)
		}
	},
}

var meUserCmd = &cobra.Command{
	Use:   "me",
	Short: "Show current user info",
	Long:  `Show information about the current user (API token owner)`,
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

		if profile.APIKey == "" {
			fmt.Printf("API key not configured for profile '%s'. Please run 'redmine auth token add <token>' or 'redmine profile add'\n", profile.Name)
			return
		}

		if profile.RedmineURL == "" {
			fmt.Printf("Redmine URL not configured for profile '%s'. Please run 'redmine profile add'\n", profile.Name)
			return
		}

		c := client.NewClient(profile.RedmineURL, profile.APIKey)

		response, err := c.GetCurrentUser()
		if err != nil {
			fmt.Printf("Error getting current user: %v\n", err)
			return
		}

		user := response.User

		fmt.Printf("Current User Information\n")
		fmt.Println(strings.Repeat("=", 30))
		fmt.Printf("ID: %d\n", user.ID)
		fmt.Printf("Name: %s\n", user.Name)
		fmt.Printf("Login: %s\n", user.Login)
		fmt.Printf("Email: %s\n", user.Email)
		fmt.Printf("Admin: %t\n", user.Admin)
		fmt.Printf("Status: %s\n", getStatusName(user.Status))
		fmt.Printf("Created: %s\n", user.CreatedOn.Format("2006-01-02 15:04:05"))
		fmt.Printf("Last login: %s\n", user.LastLoginOn.Format("2006-01-02 15:04:05"))
	},
}

func getStatusName(status int) string {
	switch status {
	case 1:
		return "Active"
	case 2:
		return "Registered"
	case 3:
		return "Locked"
	default:
		return "Unknown"
	}
}

func init() {
	rootCmd.AddCommand(usersCmd)
	usersCmd.AddCommand(listUsersCmd)
	usersCmd.AddCommand(meUserCmd)
}