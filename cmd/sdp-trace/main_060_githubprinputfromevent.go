package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/packet"
)

func githubPRInputFromEvent(event prFixtureEvent, source string, getenv func(string) string) packet.GitHubPREvidenceInput {
	// The GitHub event is the PR identity authority for packet construction.
	input := packet.GitHubPREvidenceInput{
		SchemaVersion:         "github-pr-evidence-input.v0",
		PR:                    githubPRFromEvent(event),
		CommitRange:           githubCommitRangeFromEvent(event),
		WorkflowRunID:         getenv("GITHUB_RUN_ID"),
		RequirePromptBoundary: true,
	}
	if source == "github-fixture" {
		// Fixtures carry a synthetic run id because the runner environment is not
		// available during local replay.
		input.WorkflowRunID = event.WorkflowRunID
	}
	// Prompt boundary evidence is required by default because route proof should
	// not be inferred from PR metadata alone.
	return input
}
