package main

func commandSurfaceAssessCommands() []commandSurfaceCmd {
	return []commandSurfaceCmd{
		{
			Name:        "assess",
			Description: "Assess evidence against a selected profile.",
			RequiredFlags: []flagMeta{
				{Name: "profile", Type: "string", Description: "Assessment profile."},
				{Name: "out", Type: "string", Description: "Output file."},
			},
			Subcommands: []string{"preview", "explain"},
			Variations: []string{
				"sdp-trace assess --profile adapter-capture --out \u003cfile\u003e --run \u003crun-dir\u003e",
				"sdp-trace assess --profile managed-harness --out \u003cfile\u003e --contract \u003cfile\u003e --run \u003crun-dir\u003e --adapter-registry \u003cfile\u003e --managed-policy \u003cfile\u003e --managed-witness \u003cfile\u003e",
				"sdp-trace assess --profile forensic-retention --out \u003cfile\u003e --run \u003crun-dir\u003e --redaction-policy \u003cfile\u003e",
				"sdp-trace assess --profile ci-artifact-observation --out \u003cfile\u003e --artifact-manifest \u003cfile\u003e",
				"sdp-trace assess --profile authority-envelope --out \u003cfile\u003e --authority-package \u003cfile\u003e",
				"sdp-trace assess preview --profile \u003cadapter-capture|managed-harness|forensic-retention|ci-artifact-observation|authority-envelope\u003e [profile inputs]",
				"sdp-trace assess explain --assessment-result \u003cfile\u003e",
			},
			TrustNote: "Emits verifier facts; missing or stale evidence can produce cannot_verify.",
			State:     "partial",
		},
		{
			Name:        "report",
			Usage:       "sdp-trace report --out \u003cdir\u003e \u003cruns-root-or-run-dir\u003e",
			Description: "Build report from one run or run root.",
			RequiredFlags: []flagMeta{
				{Name: "out", Type: "string", Description: "Output directory."},
			},
			Positional: []positionalMeta{
				{Name: "runs-root-or-run-dir", Description: "Run directory or root.", Required: true},
			},
			TrustNote: "Packages observed data and gaps; report presence is not proof of completeness.",
			State:     "complete",
		},
		{
			Name:        "gate",
			Usage:       "sdp-trace gate --out \u003cfile\u003e \u003cruns-root-or-run-dir\u003e",
			Description: "Produce advisory/protected gate facts for a run root.",
			RequiredFlags: []flagMeta{
				{Name: "out", Type: "string", Description: "Output file."},
			},
			Positional: []positionalMeta{
				{Name: "runs-root-or-run-dir", Description: "Run directory or root.", Required: true},
			},
			Subcommands: []string{"explain", "preview"},
			TrustNote:   "Emits verifier facts and states; not a native merge/release/risk decision.",
			State:       "partial",
		},
		{
			Name:        "checkpoint",
			Description: "Checkpoint creation and verification.",
			Subcommands: []string{"create", "verify"},
			Variations: []string{
				"sdp-trace checkpoint create --run \u003crun-dir\u003e --out \u003cfile\u003e --private-key \u003cfile\u003e [--signer-id \u003cid\u003e] [--id \u003cid\u003e]",
				"sdp-trace checkpoint verify --run \u003crun-dir\u003e --checkpoint \u003cfile\u003e [--policy \u003cfile\u003e]",
			},
			TrustNote: "Extension surface; emits signed checkpoint artifacts for downstream policy consumers.",
			State:       "partial",
		},
	}
}
