package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/packet"
)

func githubCommitRangeFromEvent(event prFixtureEvent) packet.GitHubCommitRange {
	// The commit range is copied from the event so packet generation never
	// infers base/head state from the local checkout.
	return packet.GitHubCommitRange{
		Base:            event.PullRequest.Base.SHA,
		Head:            event.PullRequest.Head.SHA,
		ChangedFilesRef: event.PullRequest.DiffURL,
	}
}
