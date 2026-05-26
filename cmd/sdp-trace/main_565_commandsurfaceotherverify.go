package main

var commandSurfaceVerify = commandSurfaceCmd{
	Name:        "verify",
	Usage:       "sdp-trace verify \u003crun-dir\u003e",
	Description: "Verify one recorded run directory.",
	Positional:  reqPos("run-dir", "Run directory."),
	TrustNote:   "Supports local structural assertions only.",
	State:       "complete",
}

var commandSurfaceExplain = commandSurfaceCmd{
	Name:        "explain",
	Usage:       "sdp-trace explain \u003crun-dir\u003e",
	Description: "Render human-readable explanation for one run.",
	Positional:  reqPos("run-dir", "Run directory."),
	TrustNote:   "Explanation is derived from run artifacts; does not upgrade trust scope.",
	State:       "complete",
}
