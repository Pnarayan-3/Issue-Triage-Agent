package triage

import "github.com/Pnarayan-3/Issue-Triage-Agent/internal/ai"

func BuildLabels(result *ai.TriageResult, threshold float64) []string {
	labels := make([]string, 0, len(result.Labels)+1)

	// Add AI-generated labels
	labels = append(labels, result.Labels...)

	// Add review status label
	if RequiresHumanReview(result, threshold) {
		labels = append(labels, "review:required")
	} else {
		labels = append(labels, "review:auto")
	}

	return labels
}