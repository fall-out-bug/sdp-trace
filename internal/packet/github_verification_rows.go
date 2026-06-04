package packet

import "strings"

// Verification rows preserve a strict downgrade order: prompt-boundary blockers
// and missing check evidence produce cannot_verify; retained checks without
// replayable artifacts remain partial; only retained successful checks pass.
func githubVerificationRow(input GitHubPREvidenceInput) Row {
	classification := ClassifyPromptBoundary(input.PromptBoundary)
	if row, ok := githubVerificationCannotVerifyRow(input, classification); ok {
		return row
	}

	if !checksHaveRetainedArtifactRefs(input) {
		return githubRow("PC-VERIFICATION", StatePartial, "GitHub check evidence is retained without retained artifact binding.", []string{"github:check"}, "GitHub CI green is not verification pass without retained artifact evidence")
	}
	if !checksSucceeded(input.Checks) {
		return githubRow("PC-VERIFICATION", StatePartial, "GitHub checks include non-success conclusions.", []string{"github:check"}, "not all retained checks concluded success")
	}
	return githubVerificationPassRow(input)
}

// Pass rows cite both the GitHub check and retained artifacts. Workflow run ID
// changes summary wording only; it does not change evidence refs.
func githubVerificationPassRow(input GitHubPREvidenceInput) Row {
	refs := append([]string{"github:check"}, artifactEvidenceRefs(input)...)
	if strings.TrimSpace(input.WorkflowRunID) != "" {
		return githubRow("PC-VERIFICATION", StatePass, "GitHub check and retained artifact evidence are retained for workflow run "+input.WorkflowRunID+".", refs, "")
	}
	return githubRow("PC-VERIFICATION", StatePass, "GitHub check and retained artifact evidence are retained.", refs, "")
}

func checksSucceeded(checks []GitHubCheck) bool {
	for _, check := range checks {
		if check.Conclusion != "success" {
			return false
		}
	}
	return true
}
