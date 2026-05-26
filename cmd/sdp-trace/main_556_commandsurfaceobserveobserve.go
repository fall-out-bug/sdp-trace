package main

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
