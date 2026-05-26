package main

import "flag"

func registerFlags(fs *flag.FlagSet, cfg *runConfig) {
	fs.BoolVar(&cfg.asJSON, "json", false, "emit JSON output")
	fs.IntVar(&cfg.iterations, "n", 20, "number of iterations")
	fs.BoolVar(&cfg.list, "list", false, "list built-in benchmarks")
	fs.StringVar(&cfg.name, "bench", "", "run a single built-in benchmark by name")
	fs.BoolVar(&cfg.raw, "raw", false, "include all_ms in JSON output")
}
