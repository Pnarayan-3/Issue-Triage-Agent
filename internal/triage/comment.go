package triage

import (
	"fmt"
	"strings"

	"github.com/Pnarayan-3/Issue-Triage-Agent/internal/ai"
)

func BuildComment(result *ai.TriageResult) string {

	labels := strings.Join(result.Labels, ", ")

	return fmt.Sprintf(
	`<!-- issue-triage-agent -->

	## 🤖 Automated Issue Triage

| Field | Result |
|---|---|
| **Type** | %s |
| **Priority** | %s |
| **Severity** | %s |
| **Component** | %s |
| **Team** | %s |
| **Confidence** | %.0f%% |

### 📝 Summary

%s

### 💡 Reason

%s

### 🏷️ Labels

%s

---

*This triage was generated automatically by Issue Triage Agent.*`,
		result.Type,
		result.Priority,
		result.Severity,
		result.Component,
		result.Team,
		result.Confidence*100,
		result.Summary,
		result.Reason,
		labels,
	)
}

const TriageCommentMarker = "<!-- issue-triage-agent -->"

func HasTriageComment(comments []string) bool {
	for _, comment := range comments {
		if strings.Contains(comment, TriageCommentMarker) {
			return true
		}
	}

	return false
}