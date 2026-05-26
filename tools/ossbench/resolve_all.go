package main

import "fmt"

func resolveAllBuiltins() ([]benchmarkDef, func(), error) {
	if err := resolveBuiltIns(); err != nil {
		return nil, nil, fmt.Errorf("resolve built-ins: %w", err)
	}
	return builtIns, cleanupTempBinary, nil
}
