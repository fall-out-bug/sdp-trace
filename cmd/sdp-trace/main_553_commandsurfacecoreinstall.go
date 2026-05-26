package main

var commandSurfaceInstall = commandSurfaceCmd{
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
}
