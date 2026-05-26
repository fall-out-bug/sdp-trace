package main

import "fmt"

func missingToolResult(p probe) probeResult {
	return probeResult{
		Name:   p.Name,
		State:  stateNotAssessed,
		Reason: fmt.Sprintf("required tool %q not found in PATH", p.NeedsTool),
	}
}
