package triage

import (
	"testing"

	"github.com/Pnarayan-3/Issue-Triage-Agent/internal/ai"
)

func TestBuildLabels_AutoReview(t *testing.T) {

	result := &ai.TriageResult{
		Confidence: 0.90,
	}

	labels := BuildLabels(result, 0.80)

	if len(labels) != 1 {
		t.Fatalf("expected 1 label, got %d", len(labels))
	}

	if labels[0] != "review:auto" {
		t.Errorf(
			"expected review:auto, got %s",
			labels[0],
		)
	}
}

func TestBuildLabels_HumanReview(t *testing.T) {

	result := &ai.TriageResult{
		Confidence: 0.70,
	}

	labels := BuildLabels(result, 0.80)

	if len(labels) != 1 {
		t.Fatalf("expected 1 label, got %d", len(labels))
	}

	if labels[0] != "review:required" {
		t.Errorf(
			"expected review:required, got %s",
			labels[0],
		)
	}
}

func TestBuildLabels_Threshold(t *testing.T) {

	result := &ai.TriageResult{
		Confidence: 0.80,
	}

	labels := BuildLabels(result, 0.80)

	if labels[0] != "review:auto" {
		t.Errorf(
			"expected review:auto at 0.80 confidence, got %s",
			labels[0],
		)
	}
}

func TestBuildLabels_IncludesAILabels(t *testing.T) {

	result := &ai.TriageResult{
		Confidence: 0.90,
		Labels:     []string{"bug", "backend"},
	}

	labels := BuildLabels(result, 0.80)

	if len(labels) != 3 {
		t.Fatalf("expected 3 labels, got %d", len(labels))
	}

	if labels[0] != "bug" {
		t.Errorf("expected bug label, got %s", labels[0])
	}

	if labels[1] != "backend" {
		t.Errorf("expected backend label, got %s", labels[1])
	}

	if labels[2] != "review:auto" {
		t.Errorf("expected review:auto, got %s", labels[2])
	}
}

func TestBuildLabels_IncludesHumanReviewLabel(t *testing.T) {

	result := &ai.TriageResult{
		Confidence: 0.60,
		Labels:     []string{"bug", "backend"},
	}

	labels := BuildLabels(result, 0.80)

	if len(labels) != 3 {
		t.Fatalf("expected 3 labels, got %d", len(labels))
	}

	if labels[2] != "review:required" {
		t.Errorf("expected review:required, got %s", labels[2])
	}
}

func TestBuildLabels_NoAILabels(t *testing.T) {

	result := &ai.TriageResult{
		Confidence: 0.90,
		Labels:     []string{},
	}

	labels := BuildLabels(result, 0.80)

	if len(labels) != 1 {
		t.Fatalf("expected 1 label, got %d", len(labels))
	}

	if labels[0] != "review:auto" {
		t.Errorf("expected review:auto, got %s", labels[0])
	}
}