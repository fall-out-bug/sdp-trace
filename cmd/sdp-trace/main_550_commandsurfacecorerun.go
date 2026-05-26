package main

var commandSurfaceRun = commandSurfaceCmd{
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
}
