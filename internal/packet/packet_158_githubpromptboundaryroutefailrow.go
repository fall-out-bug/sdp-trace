package packet

import (
	"strings"
)

func githubPromptBoundaryRouteFailRow(classification PromptBoundaryClassification) Row {
	return githubRow("PC-AGENT-ROUTE", StateFail, "Developer prompt contains recorder duties.", []string{"prompt:boundary"}, strings.Join(classification.Reasons, "; "))
}
