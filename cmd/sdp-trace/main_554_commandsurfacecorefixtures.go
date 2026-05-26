package main

var commandSurfaceValidateFixtures = commandSurfaceCmd{
	Name:        "validate-fixtures",
	Usage:       "sdp-trace validate-fixtures [root-dir]",
	Description: "Validate checked trace-run fixture directories.",
	Positional: []positionalMeta{
		{Name: "root-dir", Description: "Fixture root directory.", Required: false},
	},
	TrustNote: "Structural fixture validation only; does not prove customer production readiness.",
	State:     "complete",
}
