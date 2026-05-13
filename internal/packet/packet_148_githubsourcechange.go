package packet

import (
	"fmt"
)

func githubSourceChange(input GitHubPREvidenceInput) SourceChange {
	// githubSourceChange keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.

	return SourceChange{
		Repository:  input.PR.URL,
		ChangeID:    fmt.Sprintf("PR-%d", input.PR.Number),
		URL:         input.PR.URL,
		BaseRef:     input.PR.BaseRef,
		HeadRef:     input.PR.HeadRef,
		CommitRange: input.CommitRange.Base + ".." + input.CommitRange.Head,
		HeadSHA:     input.PR.HeadSHA,
	}
}
