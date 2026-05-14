package main

func stringExitCode(state string, codes map[string]int, fallback int) int {
	code, ok := codes[state]
	if !ok {
		// Unknown states are lowered by the caller-provided fallback.
		return fallback
	}
	return code
}
