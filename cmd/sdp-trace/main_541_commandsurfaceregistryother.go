package main

func commandSurfaceOtherCommands() []commandSurfaceCmd {
	return []commandSurfaceCmd{
		{
			Name:        "verify",
			Usage:       "sdp-trace verify \u003crun-dir\u003e",
			Description: "Verify one recorded run directory.",
			Positional: []positionalMeta{
				{Name: "run-dir", Description: "Run directory.", Required: true},
			},
			TrustNote: "Supports local structural assertions only.",
			State:     "complete",
		},
		{
			Name:        "explain",
			Usage:       "sdp-trace explain \u003crun-dir\u003e",
			Description: "Render human-readable explanation for one run.",
			Positional: []positionalMeta{
				{Name: "run-dir", Description: "Run directory.", Required: true},
			},
			TrustNote: "Explanation is derived from run artifacts; does not upgrade trust scope.",
			State:     "complete",
		},
		{
			Name:        "query",
			Usage:       "sdp-trace query --query \u003cmissing-evidence|capture-depth\u003e \u003crun-dir\u003e",
			Description: "Query missing evidence or capture depth for a run.",
			RequiredFlags: []flagMeta{
				{Name: "query", Type: "string", Description: "Query type."},
			},
			Positional: []positionalMeta{
				{Name: "run-dir", Description: "Run directory.", Required: true},
			},
			TrustNote: "Highlights gaps; missing rows are not passes.",
			State:     "complete",
		},
		{
			Name:        "query-pack",
			Usage:       "sdp-trace query-pack --pack forensics-basic-v1 --run \u003crun-dir\u003e --out \u003cfile\u003e",
			Description: "Build a forensic query package.",
			RequiredFlags: []flagMeta{
				{Name: "pack", Type: "string", Description: "Pack ID."},
				{Name: "run", Type: "string", Description: "Run directory."},
				{Name: "out", Type: "string", Description: "Output file."},
			},
			Subcommands: []string{"explain"},
			Variations: []string{
				"sdp-trace query-pack explain --result \u003cfile\u003e",
			},
			TrustNote: "Produces investigation package; digest-only or redacted data limits reconstruction.",
			State:     "partial",
		},
		{
			Name:        "witness",
			Usage:       "sdp-trace witness --kind \u003cgithub-actions|gitlab-ci|buildkite|customer-pki\u003e --out \u003cfile\u003e [--report-dir \u003cdir\u003e] [--witness-envelope \u003cfile\u003e] [--customer-pki-authority-policy \u003cfile\u003e] [--customer-pki-public-cert \u003cfile\u003e | --customer-pki-public-key \u003cfile\u003e] [--customer-pki-payload-digest \u003csha256\u003e] [--customer-pki-freshness-evidence \u003cfile\u003e] \u003cruns-root-or-run-dir\u003e",
			Description: "Bind report/run evidence to CI or customer-PKI identity.",
			RequiredFlags: []flagMeta{
				{Name: "kind", Type: "string", Description: "Witness kind."},
				{Name: "out", Type: "string", Description: "Output file."},
			},
			OptionalFlags: []flagMeta{
				{Name: "report-dir", Type: "string", Description: "Report directory."},
				{Name: "witness-envelope", Type: "string", Description: "Witness envelope file."},
				{Name: "customer-pki-authority-policy", Type: "string", Description: "Customer PKI policy."},
				{Name: "customer-pki-public-cert", Type: "string", Description: "Customer PKI public cert."},
				{Name: "customer-pki-public-key", Type: "string", Description: "Customer PKI public key."},
				{Name: "customer-pki-payload-digest", Type: "string", Description: "Payload digest."},
				{Name: "customer-pki-freshness-evidence", Type: "string", Description: "Freshness evidence file."},
			},
			Positional: []positionalMeta{
				{Name: "runs-root-or-run-dir", Description: "Run directory or root.", Required: true},
			},
			TrustNote: "CI-bound evidence is not external production trust by itself.",
			State:     "partial",
		},
		{
			Name:        "release-proof",
			Usage:       "sdp-trace release-proof --manifest \u003cfile\u003e --out \u003cfile\u003e",
			Description: "Verify a source-bound local release manifest and proof artifact.",
			RequiredFlags: []flagMeta{
				{Name: "manifest", Type: "string", Description: "Manifest file."},
				{Name: "out", Type: "string", Description: "Output file."},
			},
			TrustNote: "source_bound_local_release only; dirty/stale source or manifest mismatch fails.",
			State:     "complete",
		},
		{
			Name:        "override",
			Description: "Record an override request event.",
			Subcommands: []string{"request"},
			Variations: []string{
				"sdp-trace override request --out \u003cfile\u003e --id \u003cid\u003e --by \u003cactor\u003e --reason \u003creason\u003e --source-ref \u003cref\u003e --scope \u003cscope\u003e [--external-reference \u003cref\u003e]",
			},
			TrustNote: "Experimental surface; records an advisory override event only. It does not approve the override or change downstream policy state.",
			State:       "not_assessed",
		},
		{
			Name:        "export",
			Description: "Export cross-repo posture or telemetry.",
			Subcommands: []string{"cross-repo-posture", "cross-repo-posture explain", "telemetry"},
			Variations: []string{
				"sdp-trace export cross-repo-posture --profile cross-repo-evidence-posture-v1 --selection \u003cfile\u003e --out \u003cfile\u003e",
				"sdp-trace export cross-repo-posture explain --result \u003cfile\u003e",
				"sdp-trace export telemetry --profile prometheus-text-v1 --cross-repo-posture \u003cfile\u003e --out \u003cfile|-\u003e",
			},
			State: "partial",
		},
	}
}
