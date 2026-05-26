package main

func missingCommandResult(def benchmarkDef, iterations int) benchmarkResult {
	return benchmarkResult{
		Name:                def.Name,
		Description:         def.Description,
		Command:             "",
		BinaryPath:          "",
		BinarySource:        def.Source,
		Environment:         getEnv(),
		AttemptedIterations: iterations,
		Error:               "no command specified",
	}
}
