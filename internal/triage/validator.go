package triage

import (
	"fmt"
	"strings"

	"github.com/Pnarayan-3/Issue-Triage-Agent/internal/ai"
)

func Validate(result *ai.TriageResult) error {

	if result == nil {
		return fmt.Errorf("triage result is nil")
	}

	if !isValidType(result.Type) {
		return fmt.Errorf("invalid issue type: %s", result.Type)
	}

	if !isValidPriority(result.Priority) {
		return fmt.Errorf("invalid priority: %s", result.Priority)
	}

	if !isValidSeverity(result.Severity) {
		return fmt.Errorf("invalid severity: %s", result.Severity)
	}

	if !isValidComponent(result.Component) {
		return fmt.Errorf("invalid component: %s", result.Component)
	}

	if !isValidTeam(result.Team) {
		return fmt.Errorf("invalid team: %s", result.Team)
	}

	if result.Confidence < 0 || result.Confidence > 1 {
		return fmt.Errorf(
			"confidence must be between 0 and 1: %.2f",
			result.Confidence,
		)
	}

	if strings.TrimSpace(result.Summary) == "" {
		return fmt.Errorf("summary cannot be empty")
	}

	if strings.TrimSpace(result.Reason) == "" {
		return fmt.Errorf("reason cannot be empty")
	}

	if len(result.Labels) == 0 {
	return fmt.Errorf("at least one label is required")
	}

	if len(result.Labels) > 10 {
	return fmt.Errorf("too many labels returned: %d", len(result.Labels))
	}

	for _, label := range result.Labels {

		label = strings.TrimSpace(label)

		if label == "" {
			return fmt.Errorf("label cannot be empty")
		}

		if len(label) > 50 {
			return fmt.Errorf("label is too long: %s", label)
		}
	}

	return nil
}

func isValidType(value string) bool {

	switch value {
	case "BUG",
		"FEATURE",
		"TASK",
		"DOCUMENTATION",
		"QUESTION",
		"SECURITY":
		return true
	default:
		return false
	}
}

func isValidPriority(value string) bool {

	switch value {
	case "LOW",
		"MEDIUM",
		"HIGH",
		"CRITICAL":
		return true
	default:
		return false
	}
}

func isValidSeverity(value string) bool {

	switch value {
	case "LOW",
		"MEDIUM",
		"HIGH",
		"CRITICAL":
		return true
	default:
		return false
	}
}

func isValidComponent(value string) bool{

	switch value {
	case "BACKEND",
		"FRONTEND",
		"DATABASE",
		"DEVOPS",
		"INFRASTRUCTURE",
		"SECURITY",
		"DOCUMENTATION",
		"UNKNOWN":
		return true
	default:
		return false
	}
}

func RequiresHumanReview(result *ai.TriageResult, threshold float64) bool {
	return result.Confidence < threshold
}

func isValidTeam(value string) bool {
	switch value {
	case "BACKEND",
		"FRONTEND",
		"DATABASE",
		"DEVOPS",
		"INFRASTRUCTURE",
		"SECURITY",
		"DOCUMENTATION",
		"UNKNOWN":
		return true

	default:
		return false
	}
}