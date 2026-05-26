package main

var commandSurfaceExport = commandSurfaceCmd{
	Name:        "export",
	Description: "Export cross-repo posture or telemetry.",
	Subcommands: subs("cross-repo-posture", "cross-repo-posture explain", "telemetry"),
	Variations: vars(
		"sdp-trace export cross-repo-posture --profile cross-repo-evidence-posture-v1 --selection \u003cfile\u003e --out \u003cfile\u003e",
		"sdp-trace export cross-repo-posture explain --result \u003cfile\u003e",
		"sdp-trace export telemetry --profile prometheus-text-v1 --cross-repo-posture \u003cfile\u003e --out \u003cfile|-\u003e",
	),
	State: "partial",
}
