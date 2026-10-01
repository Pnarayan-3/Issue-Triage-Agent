package adapter

import (
	"fmt"

	"github.com/Pnarayan-3/Issue-Triage-Agent/config"
	"github.com/Pnarayan-3/Issue-Triage-Agent/internal/github"
	"github.com/Pnarayan-3/Issue-Triage-Agent/internal/issue"
)

func NewAdapter(
	cfg *config.Config,
	githubClient *github.Client,
) (Adapter, error) {

	switch cfg.IssueSource {

	case issue.SourceGitHub:
		if githubClient == nil {
			return nil, fmt.Errorf("GitHub client is required")
		}

		return NewGitHubAdapter(githubClient), nil

	case issue.SourceJira:
		return NewJiraAdapter(""), nil

	default:
		return nil, fmt.Errorf(
			"unsupported issue source: %s",
			cfg.IssueSource,
		)
	}
}