package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/UNILORN/redmine-cli/client"
	"github.com/UNILORN/redmine-cli/config"

	"github.com/spf13/cobra"
)

var showIssueCmd = &cobra.Command{
	Use:   "show [issue_id]",
	Short: "Show issue details",
	Long:  `Show detailed information about a specific issue`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		issueID, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Printf("Invalid issue ID: %s\n", args[0])
			return
		}

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

		// Check if comments flag is set
		includeComments, _ := cmd.Flags().GetBool("comments")
		var response *client.IssueResponse

		if includeComments {
			response, err = c.GetIssue(issueID, "journals")
		} else {
			response, err = c.GetIssue(issueID)
		}

		if err != nil {
			fmt.Printf("Error getting issue: %v\n", err)
			return
		}

		issue := response.Issue

		fmt.Printf("Issue #%d\n", issue.ID)
		fmt.Println(strings.Repeat("=", 50))
		fmt.Printf("Subject: %s\n", issue.Subject)
		fmt.Printf("Project: %s\n", issue.Project.Name)
		fmt.Printf("Tracker: %s\n", issue.Tracker.Name)
		fmt.Printf("Status: %s\n", issue.Status.Name)
		fmt.Printf("Priority: %s\n", issue.Priority.Name)
		fmt.Printf("Author: %s\n", issue.Author.Name)

		if issue.AssignedTo != nil {
			fmt.Printf("Assigned to: %s\n", issue.AssignedTo.Name)
		} else {
			fmt.Println("Assigned to: Not assigned")
		}

		if issue.StartDate != nil {
			fmt.Printf("Start date: %s\n", *issue.StartDate)
		}

		if issue.DueDate != nil {
			fmt.Printf("Due date: %s\n", *issue.DueDate)
		}

		fmt.Printf("Done ratio: %d%%\n", issue.DoneRatio)
		fmt.Printf("Created: %s\n", issue.CreatedOn.Format("2006-01-02 15:04:05"))
		fmt.Printf("Updated: %s\n", issue.UpdatedOn.Format("2006-01-02 15:04:05"))

		if issue.Description != "" {
			fmt.Println("\nDescription:")
			fmt.Println(strings.Repeat("-", 20))
			fmt.Println(issue.Description)
		}

		if len(issue.CustomFields) > 0 {
			fmt.Println("\nCustom Fields:")
			fmt.Println(strings.Repeat("-", 20))
			for _, field := range issue.CustomFields {
				if field.Value != "" {
					fmt.Printf("%s: %s\n", field.Name, field.Value)
				}
			}
		}

		// Display comments if requested
		if includeComments && len(issue.Journals) > 0 {
			fmt.Println("\nComments:")
			fmt.Println(strings.Repeat("=", 50))

			for _, journal := range issue.Journals {
				fmt.Printf("\n[#%d] %s - %s\n",
					journal.ID,
					journal.User.Name,
					journal.CreatedOn.Format("2006-01-02 15:04:05"))

				// Show field changes
				if len(journal.Details) > 0 {
					fmt.Println("Changes:")
					for _, detail := range journal.Details {
						oldValue := detail.OldValue
						newValue := detail.NewValue
						if oldValue == "" {
							oldValue = "(empty)"
						}
						if newValue == "" {
							newValue = "(empty)"
						}
						fmt.Printf("  - %s: %s -> %s\n", detail.Name, oldValue, newValue)
					}
				}

				// Show comment text
				if journal.Notes != "" {
					if len(journal.Details) > 0 {
						fmt.Println("Comment:")
					}
					fmt.Printf("%s\n", journal.Notes)
				}

				fmt.Println(strings.Repeat("-", 30))
			}
		}
	},
}
