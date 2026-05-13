package packet

import (
	"strings"
)

func githubVerificationPassRow(input GitHubPREvidenceInput) Row {
	// githubVerificationPassRow keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.
	refs := append([]string{"github:check"}, artifactEvidenceRefs(input)...)
	if strings.TrimSpace(input.WorkflowRunID) != "" {

		return githubRow("PC-VERIFICATION", StatePass, "GitHub check and retained artifact evidence are retained for workflow run "+input.WorkflowRunID+".", refs, "")
	}
	return githubRow("PC-VERIFICATION", StatePass, "GitHub check and retained artifact evidence are retained.", refs, "")
}
