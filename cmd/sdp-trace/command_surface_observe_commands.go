package main

var commandSurfaceInteraction = commandSurfaceCmd{
	Name:        "interaction",
	Description: "Interaction recording and summarization.",
	Subcommands: []string{"relay", "import-transcript", "summarize"},
	Variations: []string{
		"sdp-trace interaction relay --task-id \u003csafe-id\u003e --event-type \u003ctype\u003e --out \u003cfile\u003e -- \u003cforward-command...\u003e",
		"sdp-trace interaction import-transcript --source preclassified-transcript-import --task-id \u003csafe-id\u003e --events-jsonl \u003cfile\u003e --out \u003cfile\u003e",
		"sdp-trace interaction summarize --trace \u003cfile\u003e [--out \u003cfile\u003e]",
	},
	State: "partial",
}

var commandSurfaceObserve = commandSurfaceCmd{
	Name:        "observe",
	Description: "First-run harness observation.",
	Subcommands: []string{"setup", "collect", "session"},
	Variations: []string{
		"sdp-trace observe setup --profile \u003csession-profile.json\u003e --out \u003crun-dir\u003e [--command \u003charness-command-preview\u003e]",
		"sdp-trace observe collect --profile \u003csession-profile.json\u003e --run \u003crun-dir\u003e",
		"sdp-trace observe session --profile \u003csession-profile.json\u003e --out \u003crun-dir\u003e -- \u003charness-command...\u003e",
	},
	State: "partial",
}

var commandSurfaceHarness = commandSurfaceCmd{
	Name:        "harness",
	Description: "Harness event import and validation.",
	Subcommands: []string{"observe", "validate", "summarize"},
	Variations: []string{
		"sdp-trace harness observe --profile \u003charness-profile.json\u003e --source \u003charness-events.jsonl\u003e --out \u003crun-dir\u003e",
		"sdp-trace harness validate --profile \u003charness-profile.json\u003e --run \u003crun-dir\u003e --out \u003cvalidation.json\u003e",
		"sdp-trace harness summarize --validation \u003cvalidation.json\u003e",
	},
	State: "partial",
}

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

var commandSurfaceObserveGroup = []commandSurfaceCmd{
	commandSurfaceInteraction,
	commandSurfaceObserve,
	commandSurfaceHarness,
	commandSurfaceEnvelope,
}
