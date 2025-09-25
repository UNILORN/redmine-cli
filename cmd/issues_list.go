package cmd

import (
	"fmt"
	"strings"

	"github.com/UNILORN/redmine-cli/client"
	"github.com/UNILORN/redmine-cli/config"

	"github.com/spf13/cobra"
)

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

		// Check if --mine flag is set to filter by current user
		me, _ := cmd.Flags().GetBool("me")
		if me {
			// Get current user ID
			userResp, err := c.GetCurrentUser()
			if err != nil {
				fmt.Printf("Error getting current user: %v\n", err)
				return
			}
			params["assigned_to_id"] = fmt.Sprintf("%d", userResp.User.ID)
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

		// Column widths
		const (
			idWidth       = 6
			statusWidth   = 12
			assigneeWidth = 10
			dateWidth     = 12
		)

		fmt.Printf("Issues (Total: %d)\n", response.TotalCount)

		// Header
		fmt.Printf("%-*s | %-*s | %-*s | %-*s | %-*s | %-*s | %s\n",
			idWidth, "ID",
			statusWidth, "Status",
			assigneeWidth, "Assignee",
			dateWidth, "StartDate",
			dateWidth, "DueDate",
			dateWidth, "UpdatedAt",
			"Subject")

		// Separator
		fmt.Printf("%s-|-%s-|-%s-|-%s-|-%s-|-%s-|-%s\n",
			strings.Repeat("-", idWidth),
			strings.Repeat("-", statusWidth),
			strings.Repeat("-", assigneeWidth),
			strings.Repeat("-", dateWidth),
			strings.Repeat("-", dateWidth),
			strings.Repeat("-", dateWidth),
			strings.Repeat("-", 7))

		for _, issue := range response.Issues {
			assignedTo := "Not assigned"
			if issue.AssignedTo != nil {
				assignedTo = issue.AssignedTo.Name
			}

			startDate := "-"
			if issue.StartDate != nil {
				startDate = *issue.StartDate
			}

			dueDate := "-"
			if issue.DueDate != nil {
				dueDate = *issue.DueDate
			}

			// Truncate long fields to fit column widths
			status := issue.Status.Name
			if len(status) > statusWidth {
				status = status[:statusWidth-3] + "..."
			}

			fmt.Printf("#%-*d | %-*s | %-*s | %-*s | %-*s | %-*s | %s\n",
				idWidth-1, issue.ID, // -1 for the # prefix
				statusWidth, status,
				assigneeWidth, assignedTo,
				dateWidth, startDate,
				dateWidth, dueDate,
				dateWidth, issue.UpdatedOn.Format("2006-01-02"),
				issue.Subject)
		}
	},
}
