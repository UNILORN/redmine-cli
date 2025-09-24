package cmd

import (
	"fmt"
	"strconv"

	"github.com/UNILORN/redmine-cli/client"
	"github.com/UNILORN/redmine-cli/config"

	"github.com/spf13/cobra"
)

var editIssueCmd = &cobra.Command{
	Use:   "edit [issue_id]",
	Short: "Edit an existing issue",
	Long:  `Edit an existing issue in Redmine. You can update subject, description, status, assignee and add notes (comments).`,
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

		if profile.APIKey == "" {
			fmt.Printf("API key not configured for profile '%s'. Please run 'redmine auth token add <token>' or 'redmine profile add'\n", profile.Name)
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

		c := client.NewClient(profile.RedmineURL, profile.APIKey)

		// Check if issue exists
		_, err = c.GetIssue(issueID)
		if err != nil {
			fmt.Printf("Error getting issue %d: %v\n", issueID, err)
			return
		}

		// Build update request based on provided flags
		updateData := client.UpdateIssueData{}

		// Get flag values
		subject, _ := cmd.Flags().GetString("subject")
		description, _ := cmd.Flags().GetString("description")
		notes, _ := cmd.Flags().GetString("notes")
		statusIDStr, _ := cmd.Flags().GetString("status_id")
		assignedToIDStr, _ := cmd.Flags().GetString("assigned_to_id")

		// Set optional fields if provided
		if subject != "" {
			updateData.Subject = &subject
		}
		if description != "" {
			updateData.Description = &description
		}
		if notes != "" {
			updateData.Notes = &notes
		}
		if statusIDStr != "" {
			statusID, err := strconv.Atoi(statusIDStr)
			if err != nil {
				fmt.Printf("Invalid status_id: %s\n", statusIDStr)
				return
			}
			updateData.StatusID = &statusID
		}
		if assignedToIDStr != "" {
			assignedToID, err := strconv.Atoi(assignedToIDStr)
			if err != nil {
				fmt.Printf("Invalid assigned_to_id: %s\n", assignedToIDStr)
				return
			}
			updateData.AssignedToID = &assignedToID
		}

		// Check if any update data is provided
		if updateData.Subject == nil && updateData.Description == nil &&
			updateData.Notes == nil && updateData.StatusID == nil &&
			updateData.AssignedToID == nil {
			fmt.Println("No update data provided. Please specify at least one option to update.")
			return
		}

		// Update the issue
		updateReq := client.UpdateIssueRequest{
			Issue: updateData,
		}

		response, err := c.UpdateIssue(issueID, updateReq)
		if err != nil {
			fmt.Printf("Error updating issue: %v\n", err)
			return
		}

		issue := response.Issue
		assignedTo := "Not assigned"
		if issue.AssignedTo != nil {
			assignedTo = issue.AssignedTo.Name
		}

		startDateStr := ""
		if issue.StartDate != nil {
			startDateStr = *issue.StartDate
		}

		dueDateStr := ""
		if issue.DueDate != nil {
			dueDateStr = *issue.DueDate
		}

		fmt.Printf("Issue updated successfully: #%d | %s | %s | %s | %s | %s | %s\n",
			issue.ID,
			issue.Subject,
			startDateStr,
			dueDateStr,
			issue.Status.Name,
			issue.Project.Name,
			assignedTo)
	},
}
