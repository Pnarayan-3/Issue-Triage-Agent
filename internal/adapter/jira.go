package adapter

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"bytes"

	"github.com/Pnarayan-3/Issue-Triage-Agent/internal/issue"
)

type JiraAdapter struct {
	BaseURL string
	Email   string
	APIToken string
}

type JiraComment struct {
	ID   string      `json:"id"`
	Body interface{} `json:"body"`
}

type jiraCommentsResponse struct {
	Comments []JiraComment `json:"comments"`
}

type jiraIssueResponse struct {
	Key    string `json:"key"`
	Fields struct {
		Summary     string `json:"summary"`
		Description interface{} `json:"description"`
	} `json:"fields"`
}

func NewJiraAdapter(baseURL string) *JiraAdapter {
	return &JiraAdapter{
		BaseURL:  strings.TrimRight(baseURL, "/"),
		Email:    os.Getenv("JIRA_EMAIL"),
		APIToken: os.Getenv("JIRA_API_TOKEN"),
	}
}

func (j *JiraAdapter) GetIssue(id string) (*issue.Issue, error) {
	if j.BaseURL == "" {
		return nil, fmt.Errorf("JIRA_BASE_URL is not set")
	}

	if j.Email == "" {
		return nil, fmt.Errorf("JIRA_EMAIL is not set")
	}

	if j.APIToken == "" {
		return nil, fmt.Errorf("JIRA_API_TOKEN is not set")
	}

	url := j.BaseURL + "/rest/api/3/issue/" + id

	req, err := http.NewRequest(
		http.MethodGet,
		url,
		nil,
	)
	if err != nil {
		return nil, err
	}

	credentials := j.Email + ":" + j.APIToken
	auth := base64.StdEncoding.EncodeToString(
		[]byte(credentials),
	)

	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"Jira API returned status %d: %s",
			resp.StatusCode,
			string(responseBody),
		)
	}

	var jiraIssue jiraIssueResponse

	err = json.Unmarshal(responseBody, &jiraIssue)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to parse Jira response: %w",
			err,
		)
	}

	return &issue.Issue{
		ID:         jiraIssue.Key,
		Title:      jiraIssue.Fields.Summary,
		Body:       extractJiraDescription(jiraIssue.Fields.Description),
		Source:     issue.SourceJira,
		Repository: j.BaseURL,
	}, nil
}

func extractJiraDescription(description interface{}) string {
	if description == nil {
		return ""
	}

	data, err := json.Marshal(description)
	if err != nil {
		return ""
	}

	var document struct {
		Content []struct {
			Content []struct {
				Text string `json:"text"`
			} `json:"content"`
		} `json:"content"`
	}

	if err := json.Unmarshal(data, &document); err != nil {
		return ""
	}

	var parts []string

	for _, block := range document.Content {
		for _, content := range block.Content {
			if content.Text != "" {
				parts = append(parts, content.Text)
			}
		}
	}

	return strings.Join(parts, "\n")
}

func (j *JiraAdapter) AddComment(
	issueKey string,
	comment string,
) error {

	url := j.BaseURL + "/rest/api/3/issue/" + issueKey + "/comment"

	requestBody := map[string]interface{}{
		"body": map[string]interface{}{
			"type":    "doc",
			"version": 1,
			"content": []map[string]interface{}{
				{
					"type": "paragraph",
					"content": []map[string]interface{}{
						{
							"type": "text",
							"text": comment,
						},
					},
				},
			},
		},
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		url,
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		return err
	}

	credentials := j.Email + ":" + j.APIToken

	auth := base64.StdEncoding.EncodeToString(
		[]byte(credentials),
	)

	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf(
			"Jira comment API returned status %d: %s",
			resp.StatusCode,
			string(responseBody),
		)
	}

	return nil
}

func (j *JiraAdapter) AddLabels(
	issueKey string,
	labels []string,
) error {

	url := j.BaseURL + "/rest/api/3/issue/" + issueKey

	requestBody := map[string]interface{}{
		"update": map[string]interface{}{
			"labels": []map[string]interface{}{},
		},
	}

	labelUpdates := requestBody["update"].(map[string]interface{})
	operations := labelUpdates["labels"].([]map[string]interface{})

	for _, label := range labels {
		operations = append(
			operations,
			map[string]interface{}{
				"add": label,
			},
		)
	}

	labelUpdates["labels"] = operations

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(
		http.MethodPut,
		url,
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		return err
	}

	credentials := j.Email + ":" + j.APIToken

	auth := base64.StdEncoding.EncodeToString(
		[]byte(credentials),
	)

	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf(
			"Jira label API returned status %d: %s",
			resp.StatusCode,
			string(responseBody),
		)
	}

	return nil
}

func (j *JiraAdapter) UpdatePriority(
	issueKey string,
	priority string,
) error {

	jiraPriority := map[string]string{
		"LOW":      "Low",
		"MEDIUM":   "Medium",
		"HIGH":     "High",
		"CRITICAL": "Highest",
	}

	targetPriority, ok := jiraPriority[priority]
	if !ok {
		return fmt.Errorf(
			"unsupported AI priority: %s",
			priority,
		)
	}

	url := j.BaseURL + "/rest/api/3/issue/" + issueKey

	requestBody := map[string]interface{}{
		"fields": map[string]interface{}{
			"priority": map[string]string{
				"name": targetPriority,
			},
		},
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(
		http.MethodPut,
		url,
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		return err
	}

	credentials := j.Email + ":" + j.APIToken

	auth := base64.StdEncoding.EncodeToString(
		[]byte(credentials),
	)

	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf(
			"Jira priority API returned status %d: %s",
			resp.StatusCode,
			string(responseBody),
		)
	}

	return nil
}

func (j *JiraAdapter) GetTriageComment(issueKey string, marker string) (*JiraComment, error) {
	url := j.BaseURL + "/rest/api/3/issue/" + issueKey + "/comment"

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	credentials := j.Email + ":" + j.APIToken
	auth := base64.StdEncoding.EncodeToString([]byte(credentials))

	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"Jira comments API returned status %d: %s",
			resp.StatusCode,
			string(responseBody),
		)
	}

	var commentsResponse jiraCommentsResponse

	err = json.Unmarshal(responseBody, &commentsResponse)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to parse Jira comments response: %w",
			err,
		)
	}

	for _, comment := range commentsResponse.Comments {
		bodyBytes, err := json.Marshal(comment.Body)
		if err != nil {
			continue
		}

		if strings.Contains(string(bodyBytes), marker) {
			return &comment, nil
		}
	}

	return nil, nil
}

func (j *JiraAdapter) UpdateComment(
	issueKey string,
	commentID string,
	comment string,
) error {
	url := j.BaseURL +
		"/rest/api/3/issue/" +
		issueKey +
		"/comment/" +
		commentID

	requestBody := map[string]interface{}{
		"body": map[string]interface{}{
			"type":    "doc",
			"version": 1,
			"content": []map[string]interface{}{
				{
					"type": "paragraph",
					"content": []map[string]interface{}{
						{
							"type": "text",
							"text": comment,
						},
					},
				},
			},
		},
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(
		http.MethodPut,
		url,
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		return err
	}

	credentials := j.Email + ":" + j.APIToken
	auth := base64.StdEncoding.EncodeToString([]byte(credentials))

	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf(
			"Jira comment update API returned status %d: %s",
			resp.StatusCode,
			string(responseBody),
		)
	}

	return nil
}