package main

var commandSurfaceDoctor = commandSurfaceCmd{
	Name:        "doctor",
	Usage:       "sdp-trace doctor [--contract <file>]",
	Description: "Inspect local environment and contract prerequisites.",
	OptionalFlags: []flagMeta{
		{Name: "contract", Type: "string", Description: "Contract file."},
	},
	Variations: []string{
		"sdp-trace doctor --profile github-actions-git-hooks-v1 [--out <file>]",
	},
	TrustNote: "Emits structural readiness; offline or missing prerequisites can produce cannot_verify.",
	State:     "complete",
}
