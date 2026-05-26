package main

func findBuiltin(name string) (benchmarkDef, bool) {
	for _, b := range builtIns {
		if b.Name == name {
			return b, true
		}
	}
	return benchmarkDef{}, false
}
