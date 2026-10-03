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

func TestNilTriageResult(t *testing.T) {

	err := Validate(nil)

	if err == nil {
		t.Fatal("expected validation error for nil result")
	}
}

func TestInvalidType(t *testing.T) {

	result := validResult()
	result.Type = "INVALID"

	err := Validate(result)

	if err == nil {
		t.Fatal("expected validation error for invalid type")
	}
}

func TestInvalidPriority(t *testing.T) {

	result := validResult()
	result.Priority = "SUPER_HIGH"

	err := Validate(result)

	if err == nil {
		t.Fatal("expected validation error for invalid priority")
	}
}

func TestInvalidSeverity(t *testing.T) {

	result := validResult()
	result.Severity = "EXTREME"

	err := Validate(result)

	if err == nil {
		t.Fatal("expected validation error for invalid severity")
	}
}

func TestInvalidComponent(t *testing.T) {

	result := validResult()
	result.Component = "MOBILE"

	err := Validate(result)

	if err == nil {
		t.Fatal("expected validation error for invalid component")
	}
}

func TestInvalidTeam(t *testing.T) {

	result := validResult()
	result.Team = "MOBILE"

	err := Validate(result)

	if err == nil {
		t.Fatal("expected validation error for invalid team")
	}
}

func TestInvalidConfidence(t *testing.T) {

	result := validResult()
	result.Confidence = 1.5

	err := Validate(result)

	if err == nil {
		t.Fatal("expected validation error")
	}
}

func TestNegativeConfidence(t *testing.T) {

	result := validResult()
	result.Confidence = -0.1

	err := Validate(result)

	if err == nil {
		t.Fatal("expected validation error for negative confidence")
	}
}

func TestEmptySummary(t *testing.T) {

	result := validResult()
	result.Summary = ""

	err := Validate(result)

	if err == nil {
		t.Fatal("expected validation error for empty summary")
	}
}

func TestEmptyReason(t *testing.T) {

	result := validResult()
	result.Reason = ""

	err := Validate(result)

	if err == nil {
		t.Fatal("expected validation error for empty reason")
	}
}

func TestTooManyLabels(t *testing.T) {

	result := validResult()

	result.Labels = []string{
		"label1",
		"label2",
		"label3",
		"label4",
		"label5",
		"label6",
		"label7",
		"label8",
		"label9",
		"label10",
		"label11",
	}

	err := Validate(result)

	if err == nil {
		t.Fatal("expected validation error for too many labels")
	}
}

func TestEmptyLabel(t *testing.T) {

	result := validResult()
	result.Labels = []string{"bug", ""}

	err := Validate(result)

	if err == nil {
		t.Fatal("expected validation error for empty label")
	}
}

func TestLongLabel(t *testing.T) {

	result := validResult()

	longLabel := ""

	for i := 0; i < 51; i++ {
		longLabel += "a"
	}

	result.Labels = []string{longLabel}

	err := Validate(result)

	if err == nil {
		t.Fatal("expected validation error for long label")
	}
}

func TestNoLabels(t *testing.T) {

	result := validResult()
	result.Labels = []string{}

	err := Validate(result)

	if err == nil {
		t.Fatal("expected validation error when no labels are provided")
	}
}

func TestRequiresHumanReview(t *testing.T) {

	result := validResult()
	result.Confidence = 0.70

	if !RequiresHumanReview(result, 0.80) {
		t.Fatal("expected human review to be required")
	}
}

func TestDoesNotRequireHumanReview(t *testing.T) {

	result := validResult()
	result.Confidence = 0.90

	if RequiresHumanReview(result, 0.80) {
		t.Fatal("expected human review not to be required")
	}
}

func TestConfidenceAtThreshold(t *testing.T) {

	result := validResult()
	result.Confidence = 0.80

	if RequiresHumanReview(result, 0.80) {
		t.Fatal("expected result at threshold to not require human review")
	}
}

func validResult() *ai.TriageResult {

	return &ai.TriageResult{
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
}

