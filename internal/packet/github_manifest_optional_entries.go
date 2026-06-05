package packet

import "strings"

// Prompt-boundary evidence is retained when it is required or when any retained
// prompt material exists. Missing required evidence is represented by a
// not_retained manifest entry rather than by omitting the ref.
func githubPromptBoundaryEntries(input GitHubPREvidenceInput) []BundleEntry {
	if input.RequirePromptBoundary || strings.TrimSpace(input.PromptBoundary.Text) != "" || strings.TrimSpace(input.PromptBoundary.Digest) != "" {
		return []BundleEntry{authorityEntry(bundleEntry("prompt:boundary", "harness", promptBoundaryResolver(input.PromptBoundary), promptBoundaryRetainedForm(input.PromptBoundary)), "recorder", "recorder_owned", "sdp-trace recorder run", "external_retained_artifact", input.PromptBoundary.Digest)}
	}
	return nil
}

// PR body evidence is optional task-source context and must not produce a
// manifest ref when the input has no retained body resolver.
func githubPRBodyEntries(input GitHubPREvidenceInput) []BundleEntry {
	if input.PR.BodyRef != "" {
		return []BundleEntry{authorityEntry(bundleEntry("github:pr-body", "change_host", input.PR.BodyRef, "external_ref"), "ci_packet_builder", "ci_generated", "sdp-trace packet build-pr", "github_workflow_run", input.WorkflowRunID)}
	}
	return nil
}

// Agent route entries preserve the caller-provided digest and observed
// components so demo-first gates can distinguish self-declared refs from
// retained route observations.
func githubAgentRouteEntries(input GitHubPREvidenceInput) []BundleEntry {
	if len(input.AgentRouteRefs) > 0 {
		entry := bundleEntry("agent:route", "harness", strings.Join(input.AgentRouteRefs, ", "), "external_ref")
		if strings.TrimSpace(input.AgentRouteDigest) != "" {
			entry.Digest = input.AgentRouteDigest
		}
		entry.EvidenceKind = input.AgentRouteEvidenceKind
		entry.ObservedComponents = input.AgentRouteComponents
		entry = authorityEntry(entry, "recorder", "recorder_owned", "sdp-trace recorder run", "external_retained_artifact", input.AgentRouteDigest)
		return []BundleEntry{entry}
	}
	return nil
}

// Check and review entries aggregate resolver strings; individual artifacts are
// represented separately so verification pass rows can cite replayable reports.
func githubCheckEntries(input GitHubPREvidenceInput) []BundleEntry {
	if len(input.Checks) > 0 {
		return []BundleEntry{authorityEntry(bundleEntry("github:check", "ci", checkResolvers(input.Checks), "external_ref"), "ci_packet_builder", "ci_generated", "sdp-trace packet build-pr", "github_workflow_run", input.WorkflowRunID)}
	}
	return nil
}

func githubReviewEntries(input GitHubPREvidenceInput) []BundleEntry {
	if len(input.Reviews) > 0 {
		return []BundleEntry{authorityEntry(bundleEntry("github:review", "review", reviewResolvers(input.Reviews), "external_ref"), "ci_packet_builder", "ci_generated", "sdp-trace packet build-pr", "github_workflow_run", input.WorkflowRunID)}
	}
	return nil
}
