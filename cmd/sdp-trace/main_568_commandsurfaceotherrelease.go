package main

var commandSurfaceReleaseProof = commandSurfaceCmd{
	Name:        "release-proof",
	Usage:       "sdp-trace release-proof --manifest \u003cfile\u003e --out \u003cfile\u003e",
	Description: "Verify a source-bound local release manifest and proof artifact.",
	RequiredFlags: reqFlags(
		sf("manifest", "Manifest file."),
		sf("out", "Output file."),
	),
	TrustNote: "source_bound_local_release only; dirty/stale source or manifest mismatch fails.",
	State:     "complete",
}
