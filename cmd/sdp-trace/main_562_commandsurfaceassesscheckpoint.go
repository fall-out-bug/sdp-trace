package main

var commandSurfaceCheckpoint = commandSurfaceCmd{
	Name:        "checkpoint",
	Description: "Checkpoint creation and verification.",
	Subcommands: subs("create", "verify"),
	Variations: vars(
		"sdp-trace checkpoint create --run \u003crun-dir\u003e --out \u003cfile\u003e --private-key \u003cfile\u003e [--signer-id \u003cid\u003e] [--id \u003cid\u003e]",
		"sdp-trace checkpoint verify --run \u003crun-dir\u003e --checkpoint \u003cfile\u003e [--policy \u003cfile\u003e]",
	),
	TrustNote: "Extension surface; create signs run artifacts, verify replays existing checkpoint evidence.",
	State:     "partial",
}
