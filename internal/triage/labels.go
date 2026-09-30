package triage

import (
	"github.com/Pnarayan-3/Issue-Triage-Agent/internal/ai"
)

func BuildLabels(result *ai.TriageResult) []string {

	labels := make([]string, 0)

	if RequiresHumanReview(result) {
		labels = append(labels, "review:required")
	} else {
		labels = append(labels, "review:auto")
	}

	return labels
}