package triage

import (
	"testing"

	"github.com/Pnarayan-3/Issue-Triage-Agent/internal/ai"
)

func TestBuildLabels_AutoReview(t *testing.T) {
	result := &ai.TriageResult{
		Confidence: 0.90,
	}

	labels := BuildLabels(result)

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

	labels := BuildLabels(result)

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

	labels := BuildLabels(result)

	if labels[0] != "review:auto" {
		t.Errorf(
			"expected review:auto at 0.80 confidence, got %s",
			labels[0],
		)
	}
}