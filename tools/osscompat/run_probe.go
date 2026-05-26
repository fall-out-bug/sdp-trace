package main

// runProbe executes a single probe, checking prerequisites first.
func runProbe(p probe) probeResult {
	if missingProbeTool(p) {
		return missingToolResult(p)
	}
	state, reason := p.Run()
	return probeResult{
		Name:   p.Name,
		State:  state,
		Reason: reason,
	}
}
