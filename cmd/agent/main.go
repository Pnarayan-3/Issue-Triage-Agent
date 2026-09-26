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

	fmt.Println("\nSending issue to AI...")

	// client := ai.NewClient()

	// _, err := client.Analyze(issueTitle, issueBody)

	// if err != nil {
	// 	fmt.Println("❌ AI analysis failed:", err)
	// 	os.Exit(1)
	// }

	// fmt.Println("\n✅ AI analysis completed")
	client := ai.NewClient()

	err := client.ListModels()

	if err != nil {
		fmt.Println("❌ Could not list models:", err)
		os.Exit(1)
	}
}