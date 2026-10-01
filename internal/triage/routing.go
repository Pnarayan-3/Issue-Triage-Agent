package triage

import "github.com/Pnarayan-3/Issue-Triage-Agent/internal/ai"

func RouteTeam(result *ai.TriageResult) string {
	switch result.Team {
	case "BACKEND":
		return "Backend Team"

	case "FRONTEND":
		return "Frontend Team"

	case "DATABASE":
		return "Database Team"

	case "DEVOPS":
		return "DevOps Team"

	case "INFRASTRUCTURE":
		return "Infrastructure Team"

	case "SECURITY":
		return "Security Team"

	case "DOCUMENTATION":
		return "Documentation Team"

	default:
		return "Unassigned"
	}
}