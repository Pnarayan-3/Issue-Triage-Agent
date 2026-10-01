package main

import (
	"fmt"
	"os"

	"github.com/Pnarayan-3/Issue-Triage-Agent/config"
	"github.com/Pnarayan-3/Issue-Triage-Agent/internal/adapter"
	"github.com/Pnarayan-3/Issue-Triage-Agent/internal/ai"
	"github.com/Pnarayan-3/Issue-Triage-Agent/internal/triage"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Println("Configuration error:", err)
		os.Exit(1)
	}

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

	fmt.Println("=================================")
	fmt.Println("🤖 Jira Issue Triage")
	fmt.Println("=================================")

	// 1. Fetch Jira issue
	jiraAdapter := adapter.NewJiraAdapter(baseURL)

	currentIssue, err := jiraAdapter.GetIssue(issueKey)
	if err != nil {
		fmt.Println("Failed to fetch Jira issue:", err)
		os.Exit(1)
	}

	fmt.Println()
	fmt.Println("🎫 Issue")
	fmt.Println("ID     :", currentIssue.ID)
	fmt.Println("Source :", currentIssue.Source)
	fmt.Println("Title  :", currentIssue.Title)

	//traige one
	fmt.Println()
	fmt.Println("🔎 Checking for existing triage...")

	existingComment, err := jiraAdapter.GetTriageComment(
		currentIssue.ID,
		triage.TriageCommentMarker,
	)

	if err != nil {
		fmt.Println("Failed to check existing triage:", err)
		os.Exit(1)
	}

	if existingComment != nil {
		fmt.Println("⚠️ Issue has already been triaged.")
		fmt.Println("⏭️ Skipping duplicate triage.")
		return
	}

	fmt.Println("✅ No existing triage found")

	// 2. Send issue to AI
	client := ai.NewClient(cfg)

	fmt.Println()
	fmt.Println("🧠 Running AI analysis...")

	result, err := client.Analyze(
		currentIssue.Title,
		currentIssue.Body,
	)

	if err != nil {
		fmt.Println("AI analysis failed:", err)
		os.Exit(1)
	}

	// 3. Validate AI result
	if err := triage.Validate(result); err != nil {
		fmt.Println("Validation failed:", err)
		os.Exit(1)
	}

	//label
	err = jiraAdapter.AddLabels(
		currentIssue.ID,
		result.Labels,
	)

	fmt.Println()
	fmt.Println("⚡ Updating Jira priority...")

	err = jiraAdapter.UpdatePriority(
		currentIssue.ID,
		result.Priority,
	)

	if err != nil {
		fmt.Println("Failed to update Jira priority:", err)
		os.Exit(1)
	}

	fmt.Println("✅ Jira priority updated")

	if err != nil {
		fmt.Println("Failed to add Jira labels:", err)
		os.Exit(1)
	}

	fmt.Println("🏷️ Jira labels applied")


	// jira comment
	comment := triage.BuildComment(result)

		err = jiraAdapter.AddComment(
		currentIssue.ID,
		comment,
	)

	fmt.Println()
	fmt.Println("💬 Adding Jira triage comment...")

	if err != nil {
		fmt.Println("Failed to add Jira comment:", err)
		os.Exit(1)
	}

	fmt.Println()
	fmt.Println("✅ Jira triage comment added")

	// 4. Determine routing
	team := triage.RouteTeam(result)

	// 5. Display result
	fmt.Println()
	fmt.Println("=================================")
	fmt.Println("📋 TRIAGE RESULT")
	fmt.Println("=================================")

	fmt.Printf("Type        : %s\n", result.Type)
	fmt.Printf("Priority    : %s\n", result.Priority)
	fmt.Printf("Severity    : %s\n", result.Severity)
	fmt.Printf("Component   : %s\n", result.Component)
	fmt.Printf("Team        : %s\n", result.Team)
	fmt.Printf("Routing     : %s\n", team)

	fmt.Println()
	fmt.Println("🏷️ Labels")

	for _, label := range result.Labels {
		fmt.Printf("   • %s\n", label)
	}

	fmt.Println()
	fmt.Println("📝 Summary")
	fmt.Println(result.Summary)

	fmt.Println()
	fmt.Println("💡 Reason")
	fmt.Println(result.Reason)

	fmt.Println()
	fmt.Printf("🎯 Confidence : %.0f%%\n", result.Confidence*100)

	fmt.Println()
	fmt.Println("=================================")
	fmt.Println("✅ Jira triage completed")
	fmt.Println("=================================")
}