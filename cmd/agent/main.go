package main

import (
	"fmt"
	"os"

	"github.com/Pnarayan-3/Issue-Triage-Agent/internal/ai"
	"github.com/Pnarayan-3/Issue-Triage-Agent/internal/github"
	"github.com/Pnarayan-3/Issue-Triage-Agent/internal/triage"
	"github.com/Pnarayan-3/Issue-Triage-Agent/config"
)

func main() {

	cfg, err := config.Load()

	if err != nil {
		fmt.Println("❌ Configuration error:", err)
		os.Exit(1)
	}

	fmt.Printf(
		"⚙️ Confidence threshold: %.0f%%\n",
		cfg.ConfidenceThreshold*100,
	)

	issueNumber := os.Getenv("ISSUE_NUMBER")
	issueTitle := os.Getenv("ISSUE_TITLE")
	issueBody := os.Getenv("ISSUE_BODY")
	issueAction := os.Getenv("ISSUE_ACTION")
	repository := os.Getenv("REPOSITORY")

	fmt.Println("=================================")
	fmt.Println("🤖 Issue Triage Agent")
	fmt.Println("=================================")

	fmt.Println("Repository:", repository)
	fmt.Println("Issue Number:", issueNumber)
	fmt.Println("Issue Title:", issueTitle)
	fmt.Println("Event:", issueAction)
	githubClient := github.NewClient()

	// Check whether this issue has already been triaged
	fmt.Println()
	fmt.Println("🔍 Checking for existing triage comment...")

	triageComment, err := githubClient.GetTriageComment(
		issueNumber,
		triage.TriageCommentMarker,
	)

	if err != nil {
		fmt.Println("❌ Failed to retrieve issue comments:", err)
		os.Exit(1)
	}

	if triageComment != nil {

		if issueAction == "opened" {
			fmt.Println("⚠️ Triage comment already exists")
			fmt.Println("   Issue has already been triaged")
			fmt.Println("   Skipping AI analysis")

			fmt.Println()
			fmt.Println("==============================================")
			fmt.Println("⏭️ Triage skipped")
			fmt.Println("==============================================")

			return
		}

		if issueAction == "reopened" {
			fmt.Println("🔄 Existing triage comment found")
			fmt.Println("   Issue was reopened")
			fmt.Println("   Re-running AI analysis")
		}

	} else {
		fmt.Println("✅ No existing triage comment found")
	}

	// AI analysis
	client := ai.NewClient()

	fmt.Println()
	fmt.Println("🔍 Sending issue to AI...")

	result, err := client.Analyze(issueTitle, issueBody)

	if err != nil {
		fmt.Println("❌ AI analysis failed:", err)
		os.Exit(1)
	}

	// Validate AI result
	if err := triage.Validate(result); err != nil {
		fmt.Println("❌ Triage validation failed:", err)
		os.Exit(1)
	}

	team := triage.RouteTeam(result)

	fmt.Println("🎯 Routing decision:", team)

	fmt.Println("✅ Triage result validated")

	// Build review status label
	labels := triage.BuildLabels(
		result,
		cfg.ConfidenceThreshold,
	)

	if triage.RequiresHumanReview(result,cfg.ConfidenceThreshold,)	 
	{
		fmt.Println("⚠️ Low confidence triage detected")
		fmt.Println("   Human review is recommended")
	} else {
		fmt.Println("✅ Confidence threshold passed")
	}

	// Ensure GitHub labels exist
	fmt.Println()
	fmt.Println("🏷️ Checking GitHub labels...")

	err = githubClient.EnsureLabels(labels)

	if err != nil {
		fmt.Println("❌ Failed to ensure labels:", err)
		os.Exit(1)
	}

	fmt.Println("✅ All labels are ready")

	// Apply labels
	fmt.Println()
	fmt.Println("🏷️ Applying GitHub labels...")

	err = githubClient.AddLabels(issueNumber, labels)

	if err != nil {
		fmt.Println("❌ Failed to apply labels:", err)
		os.Exit(1)
	}

	fmt.Println("✅ Labels applied")

	// Post triage comment
	fmt.Println()
	comment := triage.BuildComment(result)

	if triageComment != nil && issueAction == "reopened" {

		fmt.Println("🔄 Updating existing triage comment...")

		err = githubClient.UpdateComment(
			triageComment.ID,
			comment,
		)

		if err != nil {
			fmt.Println("❌ Failed to update triage comment:", err)
			os.Exit(1)
		}

		fmt.Println("✅ Triage comment updated")

	} else {

		fmt.Println("💬 Posting triage comment...")

		err = githubClient.AddComment(issueNumber, comment)

		if err != nil {
			fmt.Println("❌ Failed to post triage comment:", err)
			os.Exit(1)
		}

		fmt.Println("✅ Triage comment posted")
	}

	// Final result
	fmt.Println()
	fmt.Println("==============================================")
	fmt.Println("📋 TRIAGE RESULT")
	fmt.Println("==============================================")

	fmt.Printf("Type        : %s\n", result.Type)
	fmt.Printf("Priority    : %s\n", result.Priority)
	fmt.Printf("Severity    : %s\n", result.Severity)
	fmt.Printf("Component   : %s\n", result.Component)
	fmt.Printf("Team        : %s\n", result.Team)

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
	fmt.Println("==============================================")
	fmt.Println("✅ Triage completed")
	fmt.Println("==============================================")
}