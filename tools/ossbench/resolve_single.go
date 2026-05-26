package main

import "fmt"

func resolveSingleBuiltin(name string) ([]benchmarkDef, func(), error) {
	def, ok := findBuiltin(name)
	if !ok {
		return nil, nil, fmt.Errorf("unknown benchmark: %s", name)
	}
	if err := resolveBuiltIns(); err != nil {
		return nil, nil, fmt.Errorf("resolve built-ins: %w", err)
	}
	def.Cmd = tempBinaryPath
	def.Source = "temp-build"
	return []benchmarkDef{def}, cleanupTempBinary, nil
}
