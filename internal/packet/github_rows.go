package packet

import "strings"

// GitHub rows are generated in packet-control order. Helper rows stay explicit
// so default not_assessed and partial claims keep their evidence and closure
// reasons visible.
func githubRows(input GitHubPREvidenceInput) []Row {
	return []Row{
		githubChangeRow(input),
		githubInitiatorRow(input),
		githubAgentRouteRow(input),
		githubMutationRow(input),
		githubVerificationRow(input),
		githubReviewRow(input),
		githubAuthorityRow(),
		githubTheaterRow(),
		githubAttestationRow(),
		githubDecisionRow(),
		githubResidualGapsRow(),
	}
}

// Change and mutation rows both depend on a complete commit range, but change
// can still cite PR metadata while mutation cannot cite a missing range.
func githubChangeRow(input GitHubPREvidenceInput) Row {
	if strings.TrimSpace(input.CommitRange.Base) == "" || strings.TrimSpace(input.CommitRange.Head) == "" {
		return githubRow("PC-CHANGE", StateCannotVerify, "Change-host metadata is retained but commit range is incomplete.", []string{"github:pr"}, "missing commit range base or head")
	}
	return githubRow("PC-CHANGE", StatePass, "Change-host metadata and commit range are retained.", []string{"github:pr", "git:commit-range"}, "")
}

func githubMutationRow(input GitHubPREvidenceInput) Row {
	if strings.TrimSpace(input.CommitRange.Base) == "" || strings.TrimSpace(input.CommitRange.Head) == "" {
		return githubRow("PC-MUTATION", StateCannotVerify, "Commit range is incomplete.", nil, "missing commit range base or head")
	}
	return githubRow("PC-MUTATION", StatePass, "Commit range and changed files are retained.", []string{"git:commit-range"}, "")
}

// A PR body is useful task context but still weaker than a retained issue or
// dedicated task binding, so the initiator row stays partial.
func githubInitiatorRow(input GitHubPREvidenceInput) Row {
	if input.PR.BodyRef != "" {
		return githubRow("PC-INITIATOR", StatePartial, "PR body task source is retained.", []string{"github:pr-body"}, "PR body is weaker than a dedicated issue binding")
	}
	return githubRow("PC-INITIATOR", StateNotAssessed, "No task or initiator evidence was provided.", nil, "missing PR body, issue, or retained task artifact")
}

// githubRow centralizes the generated packet row owner so helper rows do not
// drift on ownership defaults while preserving per-row evidence and reasons.
func githubRow(id, state, summary string, refs []string, reason string) Row {
	return Row{ID: id, State: state, Summary: summary, EvidenceRefs: refs, Reason: reason, Owner: "maintainer"}
}
