package main

import (
	"fmt"
	"os"

	"github.com/Pnarayan-3/Issue-Triage-Agent/internal/ai"
)

func main() {

	fmt.Println("=================================")
	fmt.Println("🤖 Issue Triage Agent")
	fmt.Println("=================================")

	client := ai.NewClient()

	err := client.ListModels()

	if err != nil {
		fmt.Println("❌ Could not list models:", err)
		os.Exit(1)
	}

	fmt.Println("=================================")
	fmt.Println("✅ Model listing completed")
	fmt.Println("=================================")
}