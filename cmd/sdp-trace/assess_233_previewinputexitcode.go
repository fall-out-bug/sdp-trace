package main

func previewInputExitCode(inputs map[string]string) int {
	for _, state := range inputs {
		if previewInputCannotVerify(state) {
			// Bad preview inputs block setup confidence without emitting a profile
			// assessment verdict.
			return exitCannotVerify
		}
	}
	return 0
}
