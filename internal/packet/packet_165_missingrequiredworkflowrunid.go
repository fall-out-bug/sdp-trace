package packet

import (
	"strings"
)

func missingRequiredWorkflowRunID(input GitHubPREvidenceInput) bool {
	return input.RequirePromptBoundary && strings.TrimSpace(input.WorkflowRunID) == ""
}
