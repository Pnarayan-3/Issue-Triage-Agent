package issue

type Issue struct {
	ID         string
	Title      string
	Body       string
	Source     string
	Repository string
}

const (
	SourceGitHub = "github"
	SourceJira   = "jira"
)