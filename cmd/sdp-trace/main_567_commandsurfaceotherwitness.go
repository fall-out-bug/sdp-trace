package main

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
