package main

import (
	"fmt"
	"os"
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

	fmt.Println("\nIssue Body:")
	fmt.Println(issueBody)

	fmt.Println("\n=================================")
	fmt.Println("Triage agent received the issue")
	fmt.Println("=================================")
}