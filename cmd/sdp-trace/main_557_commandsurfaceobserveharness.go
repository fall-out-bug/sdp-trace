package main

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
