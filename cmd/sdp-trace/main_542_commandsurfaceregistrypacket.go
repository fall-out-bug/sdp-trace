package main

func commandSurfacePacketCommands() []commandSurfaceCmd {
	return []commandSurfaceCmd{
		{
			Name:        "packet",
			Description: "Build, validate, demo-check, and render Change Evidence Packet bundles.",
			Subcommands: []string{"build-pr", "build-github", "validate", "check-demo", "render"},
			Variations: []string{
				"sdp-trace packet build-pr --source \u003cgithub-actions|github-fixture\u003e --out \u003cdir\u003e [--github-event \u003cfile\u003e] [--checks-json \u003cfile\u003e] [--artifacts-json \u003cfile\u003e] [--route-manifest \u003cfile\u003e] [--github-api-url \u003curl\u003e]",
				"sdp-trace packet build-github --github-input \u003cfile\u003e --out \u003cfile\u003e",
				"sdp-trace packet validate --bundle \u003cfile\u003e",
				"sdp-trace packet check-demo --bundle \u003cfile\u003e",
				"sdp-trace packet render --bundle \u003cfile\u003e --out \u003cfile\u003e",
			},
			State: "partial",
		},
		{
			Name:        "pr-review",
			Description: "Build, run, synthesize, validate, and summarize automated PR review evidence.",
			Subcommands: []string{"packet", "run", "synthesize", "validate", "summarize", "check"},
			Variations: []string{
				"sdp-trace pr-review packet --out \u003cdir\u003e --repo-id \u003csafe-id\u003e --change-ref \u003cpr|mr|change-id\u003e --base \u003csha\u003e --head \u003csha\u003e --diff \u003cfile\u003e [--metadata \u003cfile\u003e] [--context \u003cfile\u003e]... [--verification \u003cfile\u003e]... [--ci-state \u003cstate\u003e] [--created-by \u003cactor\u003e]",
				"sdp-trace pr-review run --packet \u003cdir\u003e --profile \u003cfile\u003e --out \u003cdir\u003e [--preview] [--work-dir \u003cdir\u003e] [--allow-external-runner \u003crunner\u003e]... [--not-assessed-reason \u003creason\u003e]",
				"sdp-trace pr-review synthesize --packet \u003cdir\u003e --runs \u003cdir\u003e --out \u003cfile\u003e",
				"sdp-trace pr-review validate --packet \u003cdir\u003e --profile \u003cfile\u003e --runs \u003cdir\u003e --ledger \u003cfile\u003e --out \u003cfile\u003e",
				"sdp-trace pr-review summarize --validation \u003cfile\u003e --ledger \u003cfile\u003e [--out \u003cfile\u003e]",
				"sdp-trace pr-review check --out \u003cdir\u003e --repo-id \u003csafe-id\u003e --change-ref \u003cpr|mr|change-id\u003e --base \u003csha\u003e --head \u003csha\u003e --diff \u003cfile\u003e --profile \u003cfile\u003e [--metadata \u003cfile\u003e] [--context \u003cfile\u003e]... [--verification \u003cfile\u003e]... [--work-dir \u003cdir\u003e] [--allow-external-runner \u003crunner\u003e]... [--not-assessed-reason \u003creason\u003e]",
			},
			State: "partial",
		},
	}
}
