package triage

import (
	"github.com/Pnarayan-3/Issue-Triage-Agent/internal/ai"
)

func BuildLabels(
	result *ai.TriageResult,
	threshold float64,
) []string {
	labels := make([]string, 0)

	if RequiresHumanReview(result, threshold) {
		return append(labels, "review:required")
	}

	return append(labels, "review:auto")
}