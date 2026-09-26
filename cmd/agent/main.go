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
	//fmt.Println("\n✅ AI analysis completed")
}