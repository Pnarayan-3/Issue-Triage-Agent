package main

import (
	"fmt"
	"os"

	"github.com/Pnarayan-3/Issue-Triage-Agent/config"
	"github.com/Pnarayan-3/Issue-Triage-Agent/internal/ai"
	"github.com/Pnarayan-3/Issue-Triage-Agent/internal/github"
	"github.com/Pnarayan-3/Issue-Triage-Agent/internal/logger"
	"github.com/Pnarayan-3/Issue-Triage-Agent/internal/triage"
	"github.com/Pnarayan-3/Issue-Triage-Agent/internal/adapter"
)

func main() {

	cfg, err := config.Load()

	if err != nil {
		logger.Error("configuration error=%v", err)
		os.Exit(1)
	}

	fmt.Printf(
		"⚙️ Confidence threshold: %.0f%%\n",
		cfg.ConfidenceThreshold*100,
	)

	issueNumber := os.Getenv("ISSUE_NUMBER")
	// issueTitle := os.Getenv("ISSUE_TITLE")
	// issueBody := os.Getenv("ISSUE_BODY")
	issueAction := os.Getenv("ISSUE_ACTION")
	repository := os.Getenv("REPOSITORY")

	fmt.Println("=================================")
	fmt.Println("🤖 Issue Triage Agent")
	fmt.Println("=================================")

	fmt.Println("Repository:", repository)
	fmt.Println("Issue Number:", issueNumber)
	// fmt.Println("Issue Title:", issueTitle)
	fmt.Println("Event:", issueAction)

	githubClient := github.NewClient()

	githubAdapter := adapter.NewGitHubAdapter(githubClient)

	currentIssue, err := githubAdapter.GetIssue(issueNumber)
	if err != nil {
		logger.Error(
			"issue=%s stage=issue_fetch error=%v",
			issueNumber,
			err,
		)
		os.Exit(1)
	}

	logger.Success(
		"issue=%s stage=issue_fetch status=loaded source=%s",
		currentIssue.ID,
		currentIssue.Source,
	)

	// Check whether this issue has already been triaged
	fmt.Println()
	logger.Info(
		"issue=%s stage=triage_check",
		issueNumber,
	)

	triageComment, err := githubClient.GetTriageComment(
		issueNumber,
		triage.TriageCommentMarker,
	)

	if err != nil {
		logger.Error(
			"issue=%s stage=triage_check error=%v",
			issueNumber,
			err,
		)
		os.Exit(1)
	}

	if triageComment != nil {

		if issueAction == "opened" {
			logger.Warn(
				"issue=%s stage=triage_check status=already_triaged",
				issueNumber,
			)

			fmt.Println("   Issue has already been triaged")
			fmt.Println("   Skipping AI analysis")

			fmt.Println()
			fmt.Println("==============================================")
			fmt.Println("⏭️ Triage skipped")
			fmt.Println("==============================================")

			return
		}

		if issueAction == "reopened" {
			logger.Info(
				"issue=%s stage=triage_check action=reopened",
				issueNumber,
			)

			fmt.Println("   Existing triage comment found")
			fmt.Println("   Re-running AI analysis")
		}

	} else {
		logger.Success(
			"issue=%s stage=triage_check status=no_existing_triage",
			issueNumber,
		)
	}

	// AI analysis
	client := ai.NewClient(cfg)

	fmt.Println()
	logger.Info(
		"issue=%s stage=ai_analysis",
		issueNumber,
	)

	result, err := client.Analyze(
		currentIssue.Title,
		currentIssue.Body,
	)

	if err != nil {
		logger.Error(
			"issue=%s stage=ai_analysis error=%v",
			issueNumber,
			err,
		)
		os.Exit(1)
	}

	// Validate AI result
	if err := triage.Validate(result); err != nil {
		logger.Error(
			"issue=%s stage=validation error=%v",
			issueNumber,
			err,
		)
		os.Exit(1)
	}

	team := triage.RouteTeam(result)

	logger.Success(
		"issue=%s stage=validation status=passed",
		issueNumber,
	)

	logger.Info(
		"issue=%s stage=routing team=%s",
		issueNumber,
		team,
	)

	// Build review status label
	labels := triage.BuildLabels(
		result,
		cfg.ConfidenceThreshold,
	)

	// Confidence and review decision
	logger.Info(
		"issue=%s confidence=%.2f threshold=%.2f",
		issueNumber,
		result.Confidence,
		cfg.ConfidenceThreshold,
	)

	if triage.RequiresHumanReview(
		result,
		cfg.ConfidenceThreshold,
	) {
		logger.Warn(
			"issue=%s stage=review decision=human_review_required",
			issueNumber,
		)
	} else {
		logger.Success(
			"issue=%s stage=review decision=auto_review",
			issueNumber,
		)
	}

	// Ensure GitHub labels exist
	fmt.Println()
	logger.Info(
		"issue=%s stage=github_labels action=ensure",
		issueNumber,
	)

	err = githubClient.EnsureLabels(labels)

	if err != nil {
		logger.Error(
			"issue=%s stage=github_labels action=ensure error=%v",
			issueNumber,
			err,
		)
		os.Exit(1)
	}

	logger.Success(
		"issue=%s stage=github_labels status=ready",
		issueNumber,
	)

	// Apply labels
	fmt.Println()
	logger.Info(
		"issue=%s stage=github_labels action=apply",
		issueNumber,
	)

	err = githubClient.AddLabels(issueNumber, labels)

	if err != nil {
		logger.Error(
			"issue=%s stage=github_labels action=apply error=%v",
			issueNumber,
			err,
		)
		os.Exit(1)
	}

	logger.Success(
		"issue=%s stage=github_labels status=applied",
		issueNumber,
	)

	// Post or update triage comment
	fmt.Println()

	comment := triage.BuildComment(result)

	if triageComment != nil && issueAction == "reopened" {

		logger.Info(
			"issue=%s stage=triage_comment action=update",
			issueNumber,
		)

		err = githubClient.UpdateComment(
			triageComment.ID,
			comment,
		)

		if err != nil {
			logger.Error(
				"issue=%s stage=triage_comment action=update error=%v",
				issueNumber,
				err,
			)
			os.Exit(1)
		}

		logger.Success(
			"issue=%s stage=triage_comment status=updated",
			issueNumber,
		)

	} else {

		logger.Info(
			"issue=%s stage=triage_comment action=create",
			issueNumber,
		)

		err = githubClient.AddComment(issueNumber, comment)

		if err != nil {
			logger.Error(
				"issue=%s stage=triage_comment action=create error=%v",
				issueNumber,
				err,
			)
			os.Exit(1)
		}

		logger.Success(
			"issue=%s stage=triage_comment status=posted",
			issueNumber,
		)
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