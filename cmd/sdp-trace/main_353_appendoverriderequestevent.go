package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/trace"
)

func appendOverrideRequestEvent(opts *flagSet) (trace.Event, error) {
	return trace.AppendRunEvent(opts.stringValue("out"), trace.EventPolicyOverrideRequested, overrideRequestPayload(opts), "sdp-trace-cli")
}
