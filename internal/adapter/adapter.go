package adapter

import "github.com/Pnarayan-3/Issue-Triage-Agent/internal/issue"

type Adapter interface {
	GetIssue(id string) (*issue.Issue, error)
}