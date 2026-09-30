package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type Client struct {
	Token      string
	Repository string
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
		return fmt.Errorf(
			"GitHub API returned status %d",
			resp.StatusCode,
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