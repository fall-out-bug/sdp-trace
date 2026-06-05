package main

var commandSurfaceVerify = commandSurfaceCmd{
	Name:        "verify",
	Usage:       "sdp-trace verify \u003crun-dir\u003e",
	Description: "Verify one recorded run directory.",
	Positional:  reqPos("run-dir", "Run directory."),
	TrustNote:   "Supports local structural assertions only.",
	State:       "complete",
}

var commandSurfaceExplain = commandSurfaceCmd{
	Name:        "explain",
	Usage:       "sdp-trace explain \u003crun-dir\u003e",
	Description: "Render human-readable explanation for one run.",
	Positional:  reqPos("run-dir", "Run directory."),
	TrustNote:   "Explanation is derived from run artifacts; does not upgrade trust scope.",
	State:       "complete",
}

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

var commandSurfaceWitness = commandSurfaceCmd{
	Name:        "witness",
	Usage:       "sdp-trace witness --kind \u003cgithub-actions|gitlab-ci|buildkite|customer-pki\u003e --out \u003cfile\u003e [--report-dir \u003cdir\u003e] [--witness-envelope \u003cfile\u003e] [--customer-pki-authority-policy \u003cfile\u003e] [--customer-pki-public-cert \u003cfile\u003e | --customer-pki-public-key \u003cfile\u003e] [--customer-pki-payload-digest \u003csha256\u003e] [--customer-pki-freshness-evidence \u003cfile\u003e] \u003cruns-root-or-run-dir\u003e",
	Description: "Bind report/run evidence to CI or customer-PKI identity.",
	RequiredFlags: reqFlags(
		sf("kind", "Witness kind."),
		sf("out", "Output file."),
	),
	OptionalFlags: optFlags(
		sf("report-dir", "Report directory."),
		sf("witness-envelope", "Witness envelope file."),
		sf("customer-pki-authority-policy", "Customer PKI policy."),
		sf("customer-pki-public-cert", "Customer PKI public cert."),
		sf("customer-pki-public-key", "Customer PKI public key."),
		sf("customer-pki-payload-digest", "Payload digest."),
		sf("customer-pki-freshness-evidence", "Freshness evidence file."),
	),
	Positional: reqPos("runs-root-or-run-dir", "Run directory or root."),
	TrustNote:  "CI-bound evidence is not external production trust by itself.",
	State:      "partial",
}

var commandSurfaceReleaseProof = commandSurfaceCmd{
	Name:        "release-proof",
	Usage:       "sdp-trace release-proof --manifest \u003cfile\u003e --out \u003cfile\u003e",
	Description: "Verify a source-bound local release manifest and proof artifact.",
	RequiredFlags: reqFlags(
		sf("manifest", "Manifest file."),
		sf("out", "Output file."),
	),
	TrustNote: "source_bound_local_release only; dirty/stale source or manifest mismatch fails.",
	State:     "complete",
}

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

var commandSurfaceOtherGroup = []commandSurfaceCmd{
	commandSurfaceVerify,
	commandSurfaceExplain,
	commandSurfaceQuery,
	commandSurfaceQueryPack,
	commandSurfaceWitness,
	commandSurfaceReleaseProof,
	commandSurfaceOverride,
	commandSurfaceExport,
}
