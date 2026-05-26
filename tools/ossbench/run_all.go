package main

func runAllBenchmarks(defs []benchmarkDef, iterations int, raw bool) ([]benchmarkResult, error) {
	results := make([]benchmarkResult, 0, len(defs))
	for i := range defs {
		d := &defs[i]
		if err := setupWrap(d); err != nil {
			return nil, err
		}
		res := runBenchmark(*d, iterations)
		res = finalizeResult(res, d.Cleanup, raw)
		results = append(results, res)
	}
	return results, nil
}
