package cmd

import (
	"fmt"
	"strings"

	"github.com/UNILORN/redmine-cli/client"
	"github.com/UNILORN/redmine-cli/config"
	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search [query]",
	Short: "Search for content in Redmine",
	Long:  `Search for issues, wiki pages, documents, and other content in Redmine`,
	Args:  cobra.MinimumNArgs(1),
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

		// Join all search args into a single query string
		query := strings.Join(args, " ")
		params["q"] = query

		// Get command line flags
		limit, _ := cmd.Flags().GetString("limit")
		if limit != "" {
			params["limit"] = limit
		}

		offset, _ := cmd.Flags().GetString("offset")
		if offset != "" {
			params["offset"] = offset
		}

		showDescription, _ := cmd.Flags().GetBool("description")

		scope, _ := cmd.Flags().GetString("scope")
		if scope != "" {
			params["scope"] = scope
		}

		allWords, _ := cmd.Flags().GetBool("all-words")
		if allWords {
			params["all_words"] = "1"
		}

		titlesOnly, _ := cmd.Flags().GetBool("titles-only")
		if titlesOnly {
			params["titles_only"] = "1"
		}

		openIssues, _ := cmd.Flags().GetBool("open-issues")
		if openIssues {
			params["open_issues"] = "1"
		}

		attachments, _ := cmd.Flags().GetBool("attachments")
		if attachments {
			params["attachments"] = "1"
		}

		// Content type filters
		if issues, _ := cmd.Flags().GetBool("issues"); issues {
			params["issues"] = "1"
		}
		if news, _ := cmd.Flags().GetBool("news"); news {
			params["news"] = "1"
		}
		if documents, _ := cmd.Flags().GetBool("documents"); documents {
			params["documents"] = "1"
		}
		if changesets, _ := cmd.Flags().GetBool("changesets"); changesets {
			params["changesets"] = "1"
		}
		if wikiPages, _ := cmd.Flags().GetBool("wiki-pages"); wikiPages {
			params["wiki_pages"] = "1"
		}
		if messages, _ := cmd.Flags().GetBool("messages"); messages {
			params["messages"] = "1"
		}
		if projects, _ := cmd.Flags().GetBool("projects"); projects {
			params["projects"] = "1"
		}

		response, err := c.Search(params)
		if err != nil {
			fmt.Printf("Error searching: %v\n", err)
			return
		}

		if len(response.Results) == 0 {
			fmt.Println("No results found.")
			return
		}

		// Column widths
		const (
			typeWidth = 12
			idWidth   = 8
			dateWidth = 12
		)

		fmt.Printf("Search Results (Total: %d, Query: %s)\n", response.TotalCount, query)

		// Header
		fmt.Printf("%-*s | %-*s | %-*s | %s\n",
			typeWidth, "Type",
			idWidth, "ID",
			dateWidth, "Date",
			"Title")

		// Separator
		fmt.Printf("%s-|-%s-|-%s-|-%s\n",
			strings.Repeat("-", typeWidth),
			strings.Repeat("-", idWidth),
			strings.Repeat("-", dateWidth),
			strings.Repeat("-", 10))

		for _, result := range response.Results {
			// Truncate type field to fit column width
			resultType := result.Type
			if len(resultType) > typeWidth {
				resultType = resultType[:typeWidth-3] + "..."
			}

			// Format date
			date := "-"
			if result.Datetime != "" {
				if len(result.Datetime) >= 10 {
					date = result.Datetime[:10] // Take just the date part
				} else {
					date = result.Datetime
				}
			}

			fmt.Printf("%-*s | #%-*d | %-*s | %s\n",
				typeWidth, resultType,
				idWidth-1, result.ID, // -1 for the # prefix
				dateWidth, date,
				result.Title)

			// Show description only if flag is set and it's not empty
			if showDescription && result.Description != "" {
				description := strings.TrimSpace(result.Description)
				// Replace newlines with spaces for better readability
				description = strings.ReplaceAll(description, "\n", " ")
				description = strings.ReplaceAll(description, "\r", " ")
				// Reduce multiple spaces to single space
				for strings.Contains(description, "  ") {
					description = strings.ReplaceAll(description, "  ", " ")
				}
				if len(description) > 100 {
					description = description[:97] + "..."
				}
				fmt.Printf("    %s\n", description)
			}
		}

		// Show pagination info
		if response.TotalCount > len(response.Results) {
			fmt.Printf("\nShowing %d-%d of %d results\n",
				response.Offset+1,
				response.Offset+len(response.Results),
				response.TotalCount)
		}
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)

	// Basic search parameters
	searchCmd.Flags().String("limit", "25", "Number of results to retrieve")
	searchCmd.Flags().String("offset", "0", "Offset for pagination")
	searchCmd.Flags().String("scope", "", "Search scope (all, my_project, subprojects)")

	// Search options
	searchCmd.Flags().Bool("all-words", false, "Match all query words")
	searchCmd.Flags().Bool("titles-only", false, "Search only in titles")
	searchCmd.Flags().Bool("open-issues", false, "Filter for open issues only")
	searchCmd.Flags().Bool("attachments", false, "Search in attachments")
	searchCmd.Flags().BoolP("description", "d", false, "Show description/body content in results")

	// Content type filters
	searchCmd.Flags().Bool("issues", false, "Search in issues")
	searchCmd.Flags().Bool("news", false, "Search in news")
	searchCmd.Flags().Bool("documents", false, "Search in documents")
	searchCmd.Flags().Bool("changesets", false, "Search in changesets")
	searchCmd.Flags().Bool("wiki-pages", false, "Search in wiki pages")
	searchCmd.Flags().Bool("messages", false, "Search in messages")
	searchCmd.Flags().Bool("projects", false, "Search in projects")
}
