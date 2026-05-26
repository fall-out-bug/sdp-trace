package main

import "io"

func runWithConfig(cfg runConfig, stdout, stderr io.Writer) int {
	if cfg.list {
		return handleList(stdout, stderr, cfg.args)
	}
	return executeBenchmarks(cfg, stdout, stderr)
}
