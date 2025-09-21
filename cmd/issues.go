package cmd

import (
	"github.com/spf13/cobra"
)

var issuesCmd = &cobra.Command{
	Use:   "issues",
	Short: "Manage Redmine issues",
	Long:  `List, view, and manage Redmine issues`,
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
	issuesCmd.AddCommand(addIssueCmd)
	issuesCmd.AddCommand(editIssueCmd)
	issuesCmd.AddCommand(urlIssueCmd)

	// Add flags to list command
	listIssuesCmd.Flags().String("limit", "25", "Number of issues to retrieve")
	listIssuesCmd.Flags().String("offset", "0", "Offset for pagination")
	listIssuesCmd.Flags().String("project", "", "Project ID to filter by")
	listIssuesCmd.Flags().String("status", "", "Status ID to filter by")
	listIssuesCmd.Flags().Bool("me", false, "Filter issues authored by current user")

	// Add flags to show command
	showIssueCmd.Flags().BoolP("comments", "c", false, "Include comments (journals) in the output")

	// Add flags to add command
	addIssueCmd.Flags().String("project", "", "Project number")
	addIssueCmd.Flags().String("tracker", "", "Tracker number")
	addIssueCmd.Flags().String("title", "", "Issue title")
	addIssueCmd.Flags().String("description", "", "Issue description")
	addIssueCmd.Flags().String("parent", "", "Parent issue ID")
	addIssueCmd.Flags().String("assignee", "", "Assignee email")
	addIssueCmd.Flags().String("start-date", "", "Start date (YYYY-MM-DD)")
	addIssueCmd.Flags().String("due-date", "", "Due date (YYYY-MM-DD)")

	// Add flags to edit command
	editIssueCmd.Flags().String("subject", "", "New issue subject/title")
	editIssueCmd.Flags().String("description", "", "New issue description")
	editIssueCmd.Flags().String("notes", "", "Add notes/comments to the issue")
	editIssueCmd.Flags().String("status_id", "", "Status ID")
	editIssueCmd.Flags().String("assigned_to_id", "", "User ID to assign the issue to")
}
