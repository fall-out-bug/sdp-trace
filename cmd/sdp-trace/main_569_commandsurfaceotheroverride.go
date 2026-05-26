package main

var commandSurfaceOverride = commandSurfaceCmd{
	Name:        "override",
	Description: "Record an override request event.",
	Subcommands: subs("request"),
	Variations: vars(
		"sdp-trace override request --out \u003cfile\u003e --id \u003cid\u003e --by \u003cactor\u003e --reason \u003creason\u003e --source-ref \u003cref\u003e --scope \u003cscope\u003e [--external-reference \u003cref\u003e]",
	),
	TrustNote: "Experimental surface; records an advisory override event only. It does not approve the override or change downstream policy state.",
	State:     "not_assessed",
}
