package main

func baseResult(def benchmarkDef, cmdDisplay string, argv []string, attempted int, lastErr string) benchmarkResult {
	return benchmarkResult{
		Name:                def.Name,
		Description:         def.Description,
		Command:             cmdDisplay,
		Argv:                argv,
		WorkingDirectory:    def.Dir,
		BinaryPath:          def.Cmd,
		BinarySource:        def.Source,
		Environment:         getEnv(),
		AttemptedIterations: attempted,
		Error:               lastErr,
	}
}
