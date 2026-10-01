package adapter

import (
	"fmt"
	"github.com/Pnarayan-3/Issue-Triage-Agent/internal/github"
	"github.com/Pnarayan-3/Issue-Triage-Agent/internal/issue"
	
)

type GitHubAdapter struct {
	Client *github.Client
}

func NewGitHubAdapter(client *github.Client) *GitHubAdapter {
	return &GitHubAdapter{
		Client: client,
	}
}

func (g *GitHubAdapter) GetIssue(id string) (*issue.Issue, error) {
	githubIssue, err := g.Client.GetIssue(id)
	if err != nil {
		return nil, err
	}

	return &issue.Issue{
		ID:         fmt.Sprintf("%d", githubIssue.ID),
		Title:      githubIssue.Title,
		Body:       githubIssue.Body,
		Source:     issue.SourceGitHub,
		Repository: g.Client.Repository,
	}, nil
}