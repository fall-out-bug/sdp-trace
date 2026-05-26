package main

var commandSurfaceGate = commandSurfaceCmd{
	Name:        "gate",
	Usage:       "sdp-trace gate --out \u003cfile\u003e \u003cruns-root-or-run-dir\u003e",
	Description: "Produce advisory/protected gate facts for a run root.",
	RequiredFlags: reqFlags(
		sf("out", "Output file."),
	),
	Positional:  reqPos("runs-root-or-run-dir", "Run directory or root."),
	Subcommands: subs("explain", "preview"),
	TrustNote:   "Emits verifier facts and states; not a native merge/release/risk decision.",
	State:       "partial",
}
