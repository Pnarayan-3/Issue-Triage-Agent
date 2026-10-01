package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type Client struct {
	Token      string
	Repository string
}

type IssueComment struct {
	ID   int64  `json:"id"`
	Body string `json:"body"`
}

type Issue struct {
	ID    int64  `json:"number"`
	Title string `json:"title"`
	Body  string `json:"body"`
}

func NewClient() *Client {
	return &Client{
		Token:      os.Getenv("GITHUB_TOKEN"),
		Repository: os.Getenv("REPOSITORY"),
	}
}

func (c *Client) request(method string, url string, body interface{}) error {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(
		method,
		url,
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		responseBody, _ := io.ReadAll(resp.Body)

		return fmt.Errorf(
			"GitHub API returned status %d: %s",
			resp.StatusCode,
			string(responseBody),
		)
	}

	return nil
}

func (c *Client) AddLabels(issueNumber string, labels []string) error {
	url := fmt.Sprintf(
		"https://api.github.com/repos/%s/issues/%s/labels",
		c.Repository,
		issueNumber,
	)

	body := map[string]interface{}{
		"labels": labels,
	}

	return c.request(
		http.MethodPost,
		url,
		body,
	)
}

func (c *Client) CreateLabel(label string) error {
	url := fmt.Sprintf(
		"https://api.github.com/repos/%s/labels",
		c.Repository,
	)

	body := map[string]string{
		"name":        label,
		"color":       "6B7280",
		"description": "Automatically created by Issue Triage Agent",
	}

	return c.request(
		http.MethodPost,
		url,
		body,
	)
}

func (c *Client) EnsureLabels(labels []string) error {
	for _, label := range labels {
		labelURL := fmt.Sprintf(
			"https://api.github.com/repos/%s/labels/%s",
			c.Repository,
			url.PathEscape(label),
		)

		req, err := http.NewRequest(
			http.MethodGet,
			labelURL,
			nil,
		)

		if err != nil {
			return err
		}

		req.Header.Set("Authorization", "Bearer "+c.Token)
		req.Header.Set("Accept", "application/vnd.github+json")
		req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}

		if resp.StatusCode == http.StatusOK {
			resp.Body.Close()

			fmt.Printf("   ✓ Label exists: %s\n", label)
			continue
		}

		if resp.StatusCode == http.StatusNotFound {
			resp.Body.Close()

			fmt.Printf("   + Creating label: %s\n", label)

			if err := c.CreateLabel(label); err != nil {
				return fmt.Errorf(
					"failed to create label %q: %w",
					label,
					err,
				)
			}

			continue
		}

		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		return fmt.Errorf(
			"failed to check label %q: GitHub returned status %d: %s",
			label,
			resp.StatusCode,
			string(body),
		)
	}

	return nil
}

func (c *Client) AddComment(issueNumber string, comment string) error {
	url := fmt.Sprintf(
		"https://api.github.com/repos/%s/issues/%s/comments",
		c.Repository,
		issueNumber,
	)

	body := map[string]string{
		"body": comment,
	}

	return c.request(
		http.MethodPost,
		url,
		body,
	)
}

func (c *Client) GetComments(issueNumber string) ([]IssueComment, error) {
	url := fmt.Sprintf(
		"https://api.github.com/repos/%s/issues/%s/comments",
		c.Repository,
		issueNumber,
	)

	req, err := http.NewRequest(
		http.MethodGet,
		url,
		nil,
	)

	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)

		return nil, fmt.Errorf(
			"GitHub API returned status %d: %s",
			resp.StatusCode,
			string(body),
		)
	}

	var comments []IssueComment

	err = json.NewDecoder(resp.Body).Decode(&comments)
	if err != nil {
		return nil, err
	}

	return comments, nil
}

func (c *Client) UpdateComment(commentID int64, comment string) error {
	url := fmt.Sprintf(
		"https://api.github.com/repos/%s/issues/comments/%d",
		c.Repository,
		commentID,
	)

	body := map[string]string{
		"body": comment,
	}

	return c.request(
		http.MethodPatch,
		url,
		body,
	)
}

func (c *Client) GetTriageComment(
	issueNumber string,
	marker string,
) (*IssueComment, error) {

	comments, err := c.GetComments(issueNumber)
	if err != nil {
		return nil, err
	}

	for _, comment := range comments {
		if strings.Contains(comment.Body, marker) {
			return &comment, nil
		}
	}

	return nil, nil
}

func (c *Client) GetIssue(issueNumber string) (*Issue, error) {
	apiURL := fmt.Sprintf(
		"https://api.github.com/repos/%s/issues/%s",
		c.Repository,
		issueNumber,
	)

	req, err := http.NewRequest(
		http.MethodGet,
		apiURL,
		nil,
	)

	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		responseBody, _ := io.ReadAll(resp.Body)

		return nil, fmt.Errorf(
			"GitHub API returned status %d: %s",
			resp.StatusCode,
			string(responseBody),
		)
	}

	var issue Issue

	err = json.NewDecoder(resp.Body).Decode(&issue)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to parse GitHub issue response: %w",
			err,
		)
	}

	return &issue, nil
}