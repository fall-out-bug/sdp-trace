package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/prreview"
)

var prReviewPacketRequiredFlags = []requiredCLIFlag{
	{"out", "pr-review packet requires --out"},
	{"repo-id", "pr-review packet requires --repo-id"},
	{"change-ref", "pr-review packet requires --change-ref"},
	{"base", "pr-review packet requires --base"},
	{"head", "pr-review packet requires --head"},
	{"diff", "pr-review packet requires --diff"},
}

var prReviewPacketStringFlags = []struct {
	name         string
	defaultValue string
}{
	{"out", ""},
	{"repo-id", ""},
	{"change-ref", ""},
	{"base", ""},
	{"head", ""},
	{"diff", ""},
	{"metadata", ""},
	{"context", ""},
	{"verification", ""},
	{"ci-state", prreview.StateNotAssessed},
	{"created-by", "sdp-trace-cli"},
}
