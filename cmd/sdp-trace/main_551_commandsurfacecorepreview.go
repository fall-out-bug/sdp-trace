package main

var commandSurfaceDryRun = commandSurfaceCmd{
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
}

var commandSurfacePreview = commandSurfaceCmd{
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
}
