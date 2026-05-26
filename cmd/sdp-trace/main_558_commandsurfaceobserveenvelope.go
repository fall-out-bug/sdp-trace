package main

var commandSurfaceEnvelope = commandSurfaceCmd{
	Name:        "envelope",
	Usage:       "sdp-trace envelope summarize --envelope \u003cfile\u003e [--out \u003cfile\u003e]",
	Description: "Summarize a delivery trace envelope.",
	Subcommands: []string{"summarize"},
	RequiredFlags: []flagMeta{
		{Name: "envelope", Type: "string", Description: "Envelope file."},
	},
	OptionalFlags: []flagMeta{
		{Name: "out", Type: "string", Description: "Output file."},
	},
	TrustNote: "Read-only over refs; reports linked and not_assessed areas.",
	State:     "complete",
}
