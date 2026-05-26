package main

var commandSurfaceWrap = commandSurfaceCmd{
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
}
