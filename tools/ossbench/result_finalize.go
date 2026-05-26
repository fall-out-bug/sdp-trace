package main

func finalizeResult(res benchmarkResult, cleanup func(), raw bool) benchmarkResult {
	if cleanup != nil {
		cleanup()
	}
	if !raw {
		res.AllMs = nil
	}
	return res
}
