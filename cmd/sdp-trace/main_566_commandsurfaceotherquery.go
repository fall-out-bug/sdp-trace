package main

var commandSurfaceQuery = commandSurfaceCmd{
	Name:        "query",
	Usage:       "sdp-trace query --query \u003cmissing-evidence|capture-depth\u003e \u003crun-dir\u003e",
	Description: "Query missing evidence or capture depth for a run.",
	RequiredFlags: reqFlags(
		sf("query", "Query type."),
	),
	Positional: reqPos("run-dir", "Run directory."),
	TrustNote:  "Highlights gaps; missing rows are not passes.",
	State:      "complete",
}

var commandSurfaceQueryPack = commandSurfaceCmd{
	Name:        "query-pack",
	Usage:       "sdp-trace query-pack --pack forensics-basic-v1 --run \u003crun-dir\u003e --out \u003cfile\u003e",
	Description: "Build a forensic query package.",
	RequiredFlags: reqFlags(
		sf("pack", "Pack ID."),
		sf("run", "Run directory."),
		sf("out", "Output file."),
	),
	Subcommands: subs("explain"),
	Variations: vars(
		"sdp-trace query-pack explain --result \u003cfile\u003e",
	),
	TrustNote: "Produces investigation package; digest-only or redacted data limits reconstruction.",
	State:     "partial",
}
