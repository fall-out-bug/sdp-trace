package verifier

import "github.com/fall_out_bug/sdp-trace/internal/trace"

func verifiedChain(runDir string, manifest trace.RunManifest, events []trace.Event) (trace.VerifierResult, *trace.IntegrityAudit, bool) {
	chainOk, chainIssue := verifyChain(events, manifestChainHead(manifest))
	if !chainOk {
		// Chain contradictions are hard fail evidence rather than cannot_verify.
		return failedChainResult(runDir, manifest.RunID, chainIssue), audit(manifest.RunID, "tampered_chain", chainIssue, "run_dir", runDir), false
	}
	return trace.VerifierResult{}, nil, true
}

func manifestChainHead(manifest trace.RunManifest) string {
	if manifest.EventChainHead != "" {
		return manifest.EventChainHead
	}
	// Older manifests may only carry final_chain_head; use it as the replay
	// binding when event_chain_head is absent.
	return manifest.FinalChainHead
}

func verifyChain(events []trace.Event, expectedHead string) (bool, string) {
	if issue := verifyChainEvents(events); issue != "" {
		// Event-level chain defects take precedence over manifest head mismatch.
		return false, issue
	}
	if issue := verifyExpectedHead(events, expectedHead); issue != "" {
		return false, issue
	}
	return true, ""
}

func verifyChainEvents(events []trace.Event) string {
	for i, event := range events {
		if issue := verifyChainEvent(events, i, event); issue != "" {
			// Return first defect to keep diagnostics deterministic.
			return issue
		}
	}
	return ""
}

func verifyExpectedHead(events []trace.Event, expectedHead string) string {
	if expectedHead != "" && events[len(events)-1].EventHash != expectedHead {
		// final_chain_head is optional, but when present it binds the retained
		// event files to the manifest head.
		return "run head does not match manifest final_chain_head"
	}
	return ""
}
