package ai

type TriageResult struct {
	Type       string   `json:"type"`
	Priority   string   `json:"priority"`
	Severity   string   `json:"severity"`
	Component  string   `json:"component"`
	Team       string   `json:"team"`
	Labels     []string `json:"labels"`
	Summary    string   `json:"summary"`
	Reason     string   `json:"reason"`
	Confidence float64  `json:"confidence"`
}