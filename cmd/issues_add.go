package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/UNILORN/redmine-cli/client"
	"github.com/UNILORN/redmine-cli/config"

	"github.com/spf13/cobra"
)

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
