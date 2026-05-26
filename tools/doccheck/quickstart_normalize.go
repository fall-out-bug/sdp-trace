package main

import "strings"

func normalizeQuickstartCommand(qs string) string {
	const prefix = "go run ./cmd/sdp-trace "
	if strings.HasPrefix(qs, prefix) {
		return "sdp-trace " + strings.TrimPrefix(qs, prefix)
	}
	return qs
}
