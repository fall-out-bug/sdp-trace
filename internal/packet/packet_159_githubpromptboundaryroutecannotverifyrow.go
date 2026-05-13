package packet

import (
	"strings"
)

func githubPromptBoundaryRouteCannotVerifyRow(classification PromptBoundaryClassification) Row {
	return githubRow("PC-AGENT-ROUTE", StateCannotVerify, "Prompt boundary evidence cannot verify developer-route independence.", []string{"prompt:boundary"}, strings.Join(classification.Reasons, "; "))
}
