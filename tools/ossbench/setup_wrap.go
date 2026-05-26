package main

func setupWrap(def *benchmarkDef) error {
	if def.Name != "sdp-trace-wrap" {
		return nil
	}
	return setupWrapTempDir(def)
}
