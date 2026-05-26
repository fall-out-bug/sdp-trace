package main

import "strings"

func resolveCustomBenchmark(args []string) ([]benchmarkDef, func(), error) {
	defs := []benchmarkDef{{
		Name:        strings.Join(args, " "),
		Description: "custom command",
		Cmd:         args[0],
		Args:        args[1:],
	}}
	return defs, nil, nil
}
