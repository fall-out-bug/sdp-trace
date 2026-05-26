package main

import (
	"fmt"
	"strings"
)

func resolveNamedBenchmark(cfg runConfig) ([]benchmarkDef, func(), error) {
	if len(cfg.args) > 0 {
		return nil, nil, fmt.Errorf("unexpected positional args with -bench: %s", strings.Join(cfg.args, " "))
	}
	return resolveSingleBuiltin(cfg.name)
}
