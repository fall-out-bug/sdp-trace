package main

var commandSurfaceCommandSurface = commandSurfaceCmd{
	Name:        "command-surface",
	Usage:       "sdp-trace command-surface",
	Description: "Emit machine-readable command surface JSON.",
	TrustNote:   "Experimental agent-discovery surface; schema_version may change.",
	State:       "complete",
}

var commandSurfaceVersion = commandSurfaceCmd{
	Name:        "version",
	Usage:       "sdp-trace version",
	Description: "Print version.",
	State:       "complete",
}
