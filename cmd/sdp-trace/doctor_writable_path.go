package main

import (
	"strings"
)

func writablePathCheck(id, path, okReason string) doctorCheck {
	if strings.TrimSpace(path) == "" {
		return emptyWritablePathCheck(id)
	}
	target, check, ok := writableProbeTarget(id, path)
	if !ok {
		// Path-shape failures are returned before probing so a file is never
		// treated as an artifact directory candidate.
		return check
	}
	return probeWritablePath(id, path, target, okReason)
}
