package main

var commandSurfacePRReview = commandSurfaceCmd{
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
}
