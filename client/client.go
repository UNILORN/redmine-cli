package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	BaseURL   string
	APIKey    string
	HTTPClient *http.Client
}

type Issue struct {
	ID          int            `json:"id"`
	Project     Project        `json:"project"`
	Tracker     Tracker        `json:"tracker"`
	Status      Status         `json:"status"`
	Priority    Priority       `json:"priority"`
	Author      User           `json:"author"`
	AssignedTo  *User          `json:"assigned_to,omitempty"`
	Subject     string         `json:"subject"`
	Description string         `json:"description"`
	StartDate   *string        `json:"start_date,omitempty"`
	DueDate     *string        `json:"due_date,omitempty"`
	DoneRatio   int            `json:"done_ratio"`
	IsPrivate   bool           `json:"is_private"`
	EstimatedHours *float64    `json:"estimated_hours,omitempty"`
	SpentHours  *float64       `json:"spent_hours,omitempty"`
	CreatedOn   time.Time      `json:"created_on"`
	UpdatedOn   time.Time      `json:"updated_on"`
	ClosedOn    *time.Time     `json:"closed_on,omitempty"`
	CustomFields []CustomField `json:"custom_fields,omitempty"`
	Journals    []Journal      `json:"journals,omitempty"`
}

type Journal struct {
	ID        int            `json:"id"`
	User      User           `json:"user"`
	Notes     string         `json:"notes"`
	CreatedOn time.Time      `json:"created_on"`
	Details   []JournalDetail `json:"details,omitempty"`
}

type JournalDetail struct {
	Property string `json:"property"`
	Name     string `json:"name"`
	OldValue string `json:"old_value"`
	NewValue string `json:"new_value"`
}

type Project struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Tracker struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Status struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Priority struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type CustomField struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Value string `json:"value"`
}

type IssuesResponse struct {
	Issues      []Issue `json:"issues"`
	TotalCount  int     `json:"total_count"`
	Offset      int     `json:"offset"`
	Limit       int     `json:"limit"`
}

type IssueResponse struct {
	Issue Issue `json:"issue"`
}

func NewClient(baseURL, apiKey string) *Client {
	return &Client{
		BaseURL: strings.TrimSuffix(baseURL, "/"),
		APIKey:  apiKey,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) makeRequest(method, endpoint string) (*http.Response, error) {
	url := c.BaseURL + endpoint

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("X-Redmine-API-Key", c.APIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	return resp, nil
}

func (c *Client) GetIssues(params map[string]string) (*IssuesResponse, error) {
	endpoint := "/issues.json"

	// Add query parameters
	if len(params) > 0 {
		paramStrings := make([]string, 0, len(params))
		for key, value := range params {
			paramStrings = append(paramStrings, fmt.Sprintf("%s=%s", key, value))
		}
		endpoint += "?" + strings.Join(paramStrings, "&")
	}

	resp, err := c.makeRequest("GET", endpoint)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var issuesResp IssuesResponse
	if err := json.Unmarshal(body, &issuesResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &issuesResp, nil
}

func (c *Client) GetIssue(id int, include ...string) (*IssueResponse, error) {
	endpoint := fmt.Sprintf("/issues/%d.json", id)

	// Add include parameter if specified
	if len(include) > 0 {
		includeParam := strings.Join(include, ",")
		endpoint += "?include=" + includeParam
	}

	resp, err := c.makeRequest("GET", endpoint)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var issueResp IssueResponse
	if err := json.Unmarshal(body, &issueResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &issueResp, nil
}