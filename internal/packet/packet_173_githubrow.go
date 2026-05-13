package packet

func githubRow(id, state, summary string, refs []string, reason string) Row {

	return Row{ID: id, State: state, Summary: summary, EvidenceRefs: refs, Reason: reason, Owner: "maintainer"}
}
