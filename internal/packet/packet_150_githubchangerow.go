package packet

import (
	"strings"
)

func githubChangeRow(input GitHubPREvidenceInput) Row {
	// githubChangeRow keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.
	if strings.TrimSpace(input.CommitRange.Base) == "" || strings.TrimSpace(input.CommitRange.Head) == "" {

		return githubRow("PC-CHANGE", StateCannotVerify, "Change-host metadata is retained but commit range is incomplete.", []string{"github:pr"}, "missing commit range base or head")
	}
	return githubRow("PC-CHANGE", StatePass, "Change-host metadata and commit range are retained.", []string{"github:pr", "git:commit-range"}, "")
}
