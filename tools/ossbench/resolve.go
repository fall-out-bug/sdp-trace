package main

func resolveBenchmarkDefs(cfg runConfig) ([]benchmarkDef, func(), error) {
	if cfg.name != "" {
		return resolveNamedBenchmark(cfg)
	}
	if len(cfg.args) > 0 {
		return resolveCustomBenchmark(cfg.args)
	}
	return resolveAllBuiltins()
}
