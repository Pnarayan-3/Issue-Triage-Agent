package triage

import "testing"

func TestHasTriageComment_Found(t *testing.T) {
	comments := []string{
		"Some normal comment",
		`<!-- issue-triage-agent -->

## 🤖 Automated Issue Triage`,
	}

	if !HasTriageComment(comments) {
		t.Error("expected triage comment to be detected")
	}
}

func TestHasTriageComment_NotFound(t *testing.T) {
	comments := []string{
		"Some normal comment",
		"Please investigate this issue.",
	}

	if HasTriageComment(comments) {
		t.Error("did not expect a triage comment")
	}
}

func TestHasTriageComment_Empty(t *testing.T) {
	comments := []string{}

	if HasTriageComment(comments) {
		t.Error("did not expect a triage comment")
	}
}

func TestHasTriageComment_MarkerInsideComment(t *testing.T) {
	comments := []string{
		"Some text before the marker",
		"<!-- issue-triage-agent -->\n\nAutomated triage result",
	}

	if !HasTriageComment(comments) {
		t.Error("expected triage marker to be detected")
	}
}

func TestHasTriageComment_MultipleComments(t *testing.T) {
	comments := []string{
		"First comment",
		"Second comment",
		"Third comment",
		`<!-- issue-triage-agent -->

## 🤖 Automated Issue Triage`,
		"Another comment",
	}

	if !HasTriageComment(comments) {
		t.Error("expected triage comment to be detected")
	}
}