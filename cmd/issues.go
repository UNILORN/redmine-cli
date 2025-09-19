package cmd

import (
	"bufio"
	"fmt"
	"os"
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

var addIssueCmd = &cobra.Command{
	Use:   "add",
	Short: "Create a new issue",
	Long:  `Create a new issue in Redmine with title, description, project, assignee, dates etc.`,
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

		// Get projects list
		projectsResp, err := c.GetProjects()
		if err != nil {
			fmt.Printf("Error getting projects: %v\n", err)
			return
		}

		// Get users list
		usersResp, err := c.GetUsers()
		if err != nil {
			fmt.Printf("Error getting users: %v\n", err)
			return
		}

		// Get trackers list
		trackersResp, err := c.GetTrackers()
		if err != nil {
			fmt.Printf("Error getting trackers: %v\n", err)
			return
		}

		reader := bufio.NewReader(os.Stdin)

		// Project selection
		var selectedProject client.Project
		projectFlag, _ := cmd.Flags().GetString("project")
		if projectFlag != "" {
			projectIndex, err := strconv.Atoi(projectFlag)
			if err != nil || projectIndex < 1 || projectIndex > len(projectsResp.Projects) {
				fmt.Printf("Invalid project number: %s (available: 1-%d)\n", projectFlag, len(projectsResp.Projects))
				return
			}
			selectedProject = projectsResp.Projects[projectIndex-1]
		} else {
			fmt.Println("Available projects:")
			for i, project := range projectsResp.Projects {
				fmt.Printf("%d. %s\n", i+1, project.Name)
			}
			fmt.Print("Select project number: ")
			projectInput, _ := reader.ReadString('\n')
			projectInput = strings.TrimSpace(projectInput)
			projectIndex, err := strconv.Atoi(projectInput)
			if err != nil || projectIndex < 1 || projectIndex > len(projectsResp.Projects) {
				fmt.Println("Invalid project selection")
				return
			}
			selectedProject = projectsResp.Projects[projectIndex-1]
		}

		// Tracker selection
		var selectedTracker client.Tracker
		trackerFlag, _ := cmd.Flags().GetString("tracker")
		if trackerFlag != "" {
			trackerIndex, err := strconv.Atoi(trackerFlag)
			if err != nil || trackerIndex < 1 || trackerIndex > len(trackersResp.Trackers) {
				fmt.Printf("Invalid tracker number: %s (available: 1-%d)\n", trackerFlag, len(trackersResp.Trackers))
				return
			}
			selectedTracker = trackersResp.Trackers[trackerIndex-1]
		} else {
			fmt.Println("Available trackers:")
			for i, tracker := range trackersResp.Trackers {
				fmt.Printf("%d. %s\n", i+1, tracker.Name)
			}
			fmt.Print("Select tracker number: ")
			trackerInput, _ := reader.ReadString('\n')
			trackerInput = strings.TrimSpace(trackerInput)
			trackerIndex, err := strconv.Atoi(trackerInput)
			if err != nil || trackerIndex < 1 || trackerIndex > len(trackersResp.Trackers) {
				fmt.Println("Invalid tracker selection")
				return
			}
			selectedTracker = trackersResp.Trackers[trackerIndex-1]
		}

		// Title input
		var title string
		titleFlag, _ := cmd.Flags().GetString("title")
		if titleFlag != "" {
			title = titleFlag
		} else {
			fmt.Print("Enter issue title: ")
			titleInput, _ := reader.ReadString('\n')
			title = strings.TrimSpace(titleInput)
		}
		if title == "" {
			fmt.Println("Title is required")
			return
		}

		// Description input
		var description string
		descriptionFlag, _ := cmd.Flags().GetString("description")
		if descriptionFlag != "" {
			description = descriptionFlag
		} else {
			fmt.Print("Enter issue description: ")
			descriptionInput, _ := reader.ReadString('\n')
			description = strings.TrimSpace(descriptionInput)
		}

		// Parent issue (optional)
		var parentIssueID int
		parentInput, _ := cmd.Flags().GetString("parent")
		if parentInput != "" {
			parentIssueID, err = strconv.Atoi(parentInput)
			if err != nil {
				fmt.Printf("Invalid parent issue ID: %s\n", parentInput)
				return
			}
		}

		// Assignee selection (optional)
		var assigneeID int
		assigneeEmail, _ := cmd.Flags().GetString("assignee")
		if assigneeEmail != "" {
			// First check if it's the current user
			currentUserResp, err := c.GetCurrentUser()
			if err == nil && currentUserResp.User.Email == assigneeEmail {
				assigneeID = currentUserResp.User.ID
			} else {
				// Search in users list
				for _, user := range usersResp.Users {
					if user.Email == assigneeEmail {
						assigneeID = user.ID
						break
					}
				}
			}
			if assigneeID == 0 {
				fmt.Printf("User with email '%s' not found\n", assigneeEmail)
				return
			}
		}

		// Get dates from flags
		startDate, _ := cmd.Flags().GetString("start-date")
		dueDate, _ := cmd.Flags().GetString("due-date")

		// Create issue request
		createReq := client.CreateIssueRequest{
			Issue: client.CreateIssueData{
				ProjectID:     selectedProject.ID,
				TrackerID:     selectedTracker.ID,
				Subject:       title,
				Description:   description,
				AssignedToID:  assigneeID,
				ParentIssueID: parentIssueID,
				StartDate:     startDate,
				DueDate:       dueDate,
			},
		}

		// Create the issue
		response, err := c.CreateIssue(createReq)
		if err != nil {
			fmt.Printf("Error creating issue: %v\n", err)
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

		fmt.Printf("Issue created successfully: #%d | %s | %s | %s | %s | %s | %s\n",
			issue.ID,
			issue.Subject,
			startDateStr,
			dueDateStr,
			issue.Status.Name,
			issue.Project.Name,
			assignedTo)
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
	issuesCmd.AddCommand(addIssueCmd)

	// Add flags to list command
	listIssuesCmd.Flags().String("limit", "25", "Number of issues to retrieve")
	listIssuesCmd.Flags().String("offset", "0", "Offset for pagination")
	listIssuesCmd.Flags().String("project", "", "Project ID to filter by")
	listIssuesCmd.Flags().String("status", "", "Status ID to filter by")

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
}