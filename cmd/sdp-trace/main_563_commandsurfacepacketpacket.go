package main

var commandSurfacePacket = commandSurfaceCmd{
	Name:        "packet",
	Description: "Build, validate, demo-check, and render Change Evidence Packet bundles.",
	Subcommands: []string{"build-pr", "build-github", "validate", "check-demo", "render"},
	Variations: []string{
		"sdp-trace packet build-pr --source \u003cgithub-actions|github-fixture\u003e --out \u003cdir\u003e [--github-event \u003cfile\u003e] [--checks-json \u003cfile\u003e] [--artifacts-json \u003cfile\u003e] [--route-manifest \u003cfile\u003e] [--github-api-url \u003curl\u003e]",
		"sdp-trace packet build-github --github-input \u003cfile\u003e --out \u003cfile\u003e",
		"sdp-trace packet validate --bundle \u003cfile\u003e",
		"sdp-trace packet check-demo --bundle \u003cfile\u003e",
		"sdp-trace packet render --bundle \u003cfile\u003e --out \u003cfile\u003e",
	},
	State: "partial",
}
