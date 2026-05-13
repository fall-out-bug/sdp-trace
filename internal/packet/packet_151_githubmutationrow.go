package packet

import (
	"strings"
)

func githubMutationRow(input GitHubPREvidenceInput) Row {
	// githubMutationRow keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.
	if strings.TrimSpace(input.CommitRange.Base) == "" || strings.TrimSpace(input.CommitRange.Head) == "" {

		return githubRow("PC-MUTATION", StateCannotVerify, "Commit range is incomplete.", nil, "missing commit range base or head")
	}
	return githubRow("PC-MUTATION", StatePass, "Commit range and changed files are retained.", []string{"git:commit-range"}, "")
}
