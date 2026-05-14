package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/packet"
)

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
