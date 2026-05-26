package main

import "runtime"

// builtIns are the standard OSS tool benchmarks.
var builtIns = []benchmarkDef{
	{
		Name:        "sdp-trace-version",
		Description: "sdp-trace version command",
		Cmd:         "sdp-trace",
		Args:        []string{"version"},
		Source:      "PATH",
	},
	{
		Name:        "sdp-trace-wrap",
		Description: "sdp-trace wrap no-op",
		Cmd:         "sdp-trace",
		Args:        wrapCommandArgs(),
		Source:      "PATH",
	},
}

func wrapCommandArgs() []string {
	if runtime.GOOS == "windows" {
		return []string{"wrap", "cmd", "/c", "exit", "0"}
	}
	return []string{"wrap", "true"}
}

// builtInsOrig holds the original Cmd and Source values so resolveBuiltIns
// can mutate the global slice and cleanupTempBinary can restore it.
var builtInsOrig = make([]benchmarkDef, len(builtIns))

func init() {
	copy(builtInsOrig, builtIns)
}
