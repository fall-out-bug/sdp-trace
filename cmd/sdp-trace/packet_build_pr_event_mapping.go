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

func githubPRFromEvent(event prFixtureEvent) packet.GitHubPR {
	return packet.GitHubPR{
		// Human-facing PR fields come from the event payload that GitHub signed
		// into the runner context.
		Number:  event.PullRequest.Number,
		URL:     event.PullRequest.HTMLURL,
		Title:   event.PullRequest.Title,
		BodyRef: event.PullRequest.BodyRef,
		Author:  event.PullRequest.User.Login,
		// Branch identity stays separate from commit-range identity so packet
		// rows can explain both names and immutable SHAs.
		BaseRef: event.PullRequest.Base.Ref,
		HeadRef: event.PullRequest.Head.Ref,
		HeadSHA: event.PullRequest.Head.SHA,
	}
}

func githubCommitRangeFromEvent(event prFixtureEvent) packet.GitHubCommitRange {
	// The commit range is copied from the event so packet generation never
	// infers base/head state from the local checkout.
	return packet.GitHubCommitRange{
		Base:            event.PullRequest.Base.SHA,
		Head:            event.PullRequest.Head.SHA,
		ChangedFilesRef: event.PullRequest.DiffURL,
	}
}
