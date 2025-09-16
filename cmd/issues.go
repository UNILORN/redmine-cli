package cmd

import (
	"fmt"
	"redmine-cli/client"
	"redmine-cli/config"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var issuesCmd = &cobra.Command{
	Use:   "issues",
	Short: "Manage Redmine issues",
	Long:  `List, view, and manage Redmine issues`,
}

var listIssuesCmd = &cobra.Command{
	Use:   "list",
	Short: "List issues",
	Long:  `List all issues from Redmine`,
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

		params := make(map[string]string)

		// Get command line flags
		limit, _ := cmd.Flags().GetString("limit")
		if limit != "" {
			params["limit"] = limit
		}

		offset, _ := cmd.Flags().GetString("offset")
		if offset != "" {
			params["offset"] = offset
		}

		projectID, _ := cmd.Flags().GetString("project")
		if projectID != "" {
			params["project_id"] = projectID
		}

		status, _ := cmd.Flags().GetString("status")
		if status != "" {
			params["status_id"] = status
		}

		response, err := c.GetIssues(params)
		if err != nil {
			fmt.Printf("Error getting issues: %v\n", err)
			return
		}

		if len(response.Issues) == 0 {
			fmt.Println("No issues found.")
			return
		}

		fmt.Printf("Issues (Total: %d)\n", response.TotalCount)
		fmt.Println(strings.Repeat("-", 100))

		for _, issue := range response.Issues {
			assignedTo := "Not assigned"
			if issue.AssignedTo != nil {
				assignedTo = issue.AssignedTo.Name
			}

			fmt.Printf("#%d | %s | %s | %s | %s | %s\n",
				issue.ID,
				truncateString(issue.Subject, 40),
				issue.Status.Name,
				issue.Priority.Name,
				assignedTo,
				issue.UpdatedOn.Format("2006-01-02"))
		}
	},
}

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

func truncateString(s string, length int) string {
	if len(s) <= length {
		return s
	}
	return s[:length-3] + "..."
}

func init() {
	rootCmd.AddCommand(issuesCmd)
	issuesCmd.AddCommand(listIssuesCmd)
	issuesCmd.AddCommand(showIssueCmd)

	// Add flags to list command
	listIssuesCmd.Flags().String("limit", "25", "Number of issues to retrieve")
	listIssuesCmd.Flags().String("offset", "0", "Offset for pagination")
	listIssuesCmd.Flags().String("project", "", "Project ID to filter by")
	listIssuesCmd.Flags().String("status", "", "Status ID to filter by")

	// Add flags to show command
	showIssueCmd.Flags().BoolP("comments", "c", false, "Include comments (journals) in the output")
}