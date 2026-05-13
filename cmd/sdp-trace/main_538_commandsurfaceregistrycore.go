package main

func commandSurfaceCoreCommands() []commandSurfaceCmd {
	return []commandSurfaceCmd{
		{
			Name:        "command-surface",
			Usage:       "sdp-trace command-surface",
			Description: "Emit machine-readable command surface JSON.",
			TrustNote:   "Experimental agent-discovery surface; schema_version may change.",
			State:       "complete",
		},
		{
			Name:        "version",
			Usage:       "sdp-trace version",
			Description: "Print version.",
			State:       "complete",
		},
		{
			Name:        "wrap",
			Usage:       "sdp-trace wrap --name <name> [--contract <file>] [--output-dir <dir>] -- <command...>",
			Description: "Observe one existing command as a trace run.",
			RequiredFlags: []flagMeta{
				{Name: "name", Type: "string", Description: "Run name."},
			},
			OptionalFlags: []flagMeta{
				{Name: "contract", Type: "string", Description: "Contract file."},
				{Name: "output-dir", Type: "string", Description: "Output directory."},
			},
			RestBehavior: "required_after_double_dash",
			OutputPaths: []outputPathMeta{
				{Flag: "output-dir", Description: "Run output directory."},
			},
			TrustNote: "Writes run artifacts; local observation only unless later bound by report/witness/profile checks.",
			State:     "complete",
		},
		{
			Name:        "run",
			Usage:       "sdp-trace run --task <task-ref> [--contract <file> | --use-default-contract] -- <command...>",
			Description: "Run a task-referenced command with an optional contract.",
			RequiredFlags: []flagMeta{
				{Name: "task", Type: "string", Description: "Task reference."},
			},
			OptionalFlags: []flagMeta{
				{Name: "contract", Type: "string", Description: "Contract file."},
				{Name: "use-default-contract", Type: "bool", Description: "Use default contract."},
			},
			RestBehavior: "required_after_double_dash",
			OutputPaths: []outputPathMeta{
				{Description: "Run artifacts directory."},
			},
			TrustNote: "Writes task-linked run artifacts; missing contract evidence remains visible.",
			State:     "complete",
		},
		{
			Name:        "dry-run",
			Usage:       "sdp-trace dry-run [--contract <file> | --use-default-contract] -- <command...>",
			Description: "Show what would run without writing run artifacts.",
			OptionalFlags: []flagMeta{
				{Name: "contract", Type: "string", Description: "Contract file."},
				{Name: "use-default-contract", Type: "bool", Description: "Use default contract."},
			},
			RestBehavior: "required_after_double_dash",
			TrustNote:    "Preview only; cannot support proof closure.",
			State:        "complete",
		},
		{
			Name:        "preview",
			Usage:       "sdp-trace preview [--contract <file> | --use-default-contract] -- <command...>",
			Description: "Preview command/contract implications before execution.",
			OptionalFlags: []flagMeta{
				{Name: "contract", Type: "string", Description: "Contract file."},
				{Name: "use-default-contract", Type: "bool", Description: "Use default contract."},
			},
			RestBehavior: "required_after_double_dash",
			TrustNote:    "Read-only preview; any unavailable profile remains not_assessed.",
			State:        "complete",
		},
		{
			Name:        "doctor",
			Usage:       "sdp-trace doctor [--contract <file>]",
			Description: "Inspect local environment and contract prerequisites.",
			OptionalFlags: []flagMeta{
				{Name: "contract", Type: "string", Description: "Contract file."},
			},
			Variations: []string{
				"sdp-trace doctor --profile github-actions-git-hooks-v1 [--out <file>]",
			},
			TrustNote: "Emits structural readiness; offline or missing prerequisites can produce cannot_verify.",
			State:     "complete",
		},
		{
			Name:        "install",
			Usage:       "sdp-trace install repo-observer --profile github-actions-git-hooks-v1 [--repository-id <safe-id>] [--write] [--force] [--out <file>]",
			Description: "Install portable repo observer files for local git hooks and GitHub Actions artifact upload.",
			Subcommands: []string{"repo-observer"},
			RequiredFlags: []flagMeta{
				{Name: "profile", Type: "string", Description: "Profile ID."},
			},
			OptionalFlags: []flagMeta{
				{Name: "repository-id", Type: "string", Description: "Safe repository ID."},
				{Name: "write", Type: "bool", Description: "Write files."},
				{Name: "force", Type: "bool", Description: "Force overwrite."},
				{Name: "out", Type: "string", Description: "Output file."},
			},
			OutputPaths: []outputPathMeta{
				{Flag: "out", Description: "Status output file."},
			},
			TrustNote: "Dry-run by default. With --write, writes only the documented allowlist.",
			State:     "complete",
		},
		{
			Name:        "validate-fixtures",
			Usage:       "sdp-trace validate-fixtures [root-dir]",
			Description: "Validate checked trace-run fixture directories.",
			Positional: []positionalMeta{
				{Name: "root-dir", Description: "Fixture root directory.", Required: false},
			},
			TrustNote: "Structural fixture validation only; does not prove customer production readiness.",
			State:     "complete",
		},
	}
}
