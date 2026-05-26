package main

var commandSurfaceReport = commandSurfaceCmd{
	Name:        "report",
	Usage:       "sdp-trace report --out \u003cdir\u003e \u003cruns-root-or-run-dir\u003e",
	Description: "Build report from one run or run root.",
	RequiredFlags: reqFlags(
		sf("out", "Output directory."),
	),
	Positional: reqPos("runs-root-or-run-dir", "Run directory or root."),
	TrustNote:  "Packages observed data and gaps; report presence is not proof of completeness.",
	State:      "complete",
}
