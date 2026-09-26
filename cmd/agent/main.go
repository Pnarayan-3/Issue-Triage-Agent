package main

import (
	"fmt"
	"os"

	"github.com/Pnarayan-3/Issue-Triage-Agent/internal/ai"
)

func main() {
	issueNumber := os.Getenv("ISSUE_NUMBER")
	issueTitle := os.Getenv("ISSUE_TITLE")
	issueBody := os.Getenv("ISSUE_BODY")
	repository := os.Getenv("REPOSITORY")

	fmt.Println("=================================")
	fmt.Println("🤖 Issue Triage Agent")
	fmt.Println("=================================")

	fmt.Println("Repository:", repository)
	fmt.Println("Issue Number:", issueNumber)
	fmt.Println("Issue Title:", issueTitle)

	client := ai.NewClient()

	fmt.Println("\nSending issue to AI...")

	result, err := client.Analyze(issueTitle, issueBody)

	if err != nil {
		fmt.Println("❌ AI analysis failed:", err)
		os.Exit(1)
	}

	fmt.Println("\n=================================")
	fmt.Println("📋 TRIAGE RESULT")
	fmt.Println("=================================")

	fmt.Println("Type:", result.Type)
	fmt.Println("Priority:", result.Priority)
	fmt.Println("Severity:", result.Severity)
	fmt.Println("Component:", result.Component)
	fmt.Println("Team:", result.Team)
	fmt.Println("Labels:", result.Labels)
	fmt.Println("Summary:", result.Summary)
	fmt.Println("Reason:", result.Reason)
	fmt.Printf("Confidence: %.2f\n", result.Confidence)

	fmt.Println("\n✅ AI analysis completed")
}