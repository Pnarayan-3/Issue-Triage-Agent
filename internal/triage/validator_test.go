package triage

import (
	"testing"

	"github.com/Pnarayan-3/Issue-Triage-Agent/internal/ai"
)

func TestValidTriageResult(t *testing.T) {

	result := &ai.TriageResult{
		Type:       "BUG",
		Priority:   "HIGH",
		Severity:   "HIGH",
		Component:  "BACKEND",
		Team:       "BACKEND",
		Labels:     []string{"bug", "backend"},
		Summary:    "Login API returns HTTP 401.",
		Reason:     "Users cannot authenticate after deployment.",
		Confidence: 0.92,
	}

	err := Validate(result)

	if err != nil {
		t.Fatalf("expected valid result, got error: %v", err)
	}
}

func TestInvalidPriority(t *testing.T) {

	result := &ai.TriageResult{
		Type:       "BUG",
		Priority:   "SUPER_HIGH",
		Severity:   "HIGH",
		Component:  "BACKEND",
		Labels:     []string{"bug"},
		Summary:    "Login API returns HTTP 401.",
		Reason:     "Authentication is failing.",
		Confidence: 0.92,
	}

	err := Validate(result)

	if err == nil {
		t.Fatal("expected validation error")
	}
}

func TestInvalidConfidence(t *testing.T) {

	result := &ai.TriageResult{
		Type:       "BUG",
		Priority:   "HIGH",
		Severity:   "HIGH",
		Component:  "BACKEND",
		Labels:     []string{"bug"},
		Summary:    "Login API returns HTTP 401.",
		Reason:     "Authentication is failing.",
		Confidence: 1.5,
	}

	err := Validate(result)

	if err == nil {
		t.Fatal("expected validation error")
	}
}