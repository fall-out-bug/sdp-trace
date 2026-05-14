package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/prreview"
)

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
