package main

import (
	"fmt"
	"os"

	"github.com/Pnarayan-3/Issue-Triage-Agent/internal/adapter"
)

func main() {
	baseURL := os.Getenv("JIRA_BASE_URL")
	issueKey := os.Getenv("JIRA_ISSUE_KEY")

	if baseURL == "" {
		fmt.Println("JIRA_BASE_URL is not set")
		os.Exit(1)
	}

	if issueKey == "" {
		fmt.Println("JIRA_ISSUE_KEY is not set")
		os.Exit(1)
	}

	jiraAdapter := adapter.NewJiraAdapter(baseURL)

	result, err := jiraAdapter.GetIssue(issueKey)
	if err != nil {
		fmt.Println("Failed to fetch Jira issue:", err)
		os.Exit(1)
	}

	fmt.Println("=================================")
	fmt.Println("🎫 Jira Issue")
	fmt.Println("=================================")

	fmt.Println("ID       :", result.ID)
	fmt.Println("Source   :", result.Source)
	fmt.Println("Title    :", result.Title)
	fmt.Println("Repository:", result.Repository)

	fmt.Println()
	fmt.Println("Body:")
	fmt.Println(result.Body)

	fmt.Println("=================================")
}