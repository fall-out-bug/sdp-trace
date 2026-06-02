package main

import "time"

func overrideRequestPayload(opts *flagSet) map[string]any {
	// Override requests are appended as trace events; they request policy
	// review and never upgrade an existing gate verdict by themselves.
	// The event payload mirrors trace fields so the appended event is the
	// authoritative override request artifact.
	payload := map[string]any{
		"override_id":  opts.stringValue("id"),
		"producer":     "sdp-trace-cli",
		"origin":       "native_cli",
		"requested_by": opts.stringValue("by"),
		"reason":       opts.stringValue("reason"),
		"source_ref":   opts.stringValue("source-ref"),
		"scope":        opts.stringValue("scope"),
		"created_at":   time.Now().UTC().Format(time.RFC3339Nano),
	}
	if external := opts.stringValue("external-reference"); external != "" {
		// External references remain metadata until another verifier resolves
		// them; the CLI does not treat them as approval proof.
		payload["external_reference"] = external
	}
	return payload
}
