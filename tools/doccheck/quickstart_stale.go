package main

func staleQuickstartCommands(qsCmds, registry []string) []string {
	registrySet := stringSliceToSet(registry)
	var stale []string
	for _, qs := range qsCmds {
		if qs == "go run ./cmd/sdp-trace --help" {
			// The top-level help probe is a docs smoke test, not a registry row.
			continue
		}
		if !isKnownCommand(qs, registrySet) {
			// Anything not recognized by exact, prefix, or base-command matching
			// is stale public documentation.
			stale = append(stale, qs)
		}
	}
	return stale
}

func isKnownCommand(qs string, registrySet map[string]bool) bool {
	normalized := normalizeQuickstartCommand(qs)
	if registrySet[normalized] {
		return true
	}
	if prefixMatchesRegistry(normalized, registrySet) {
		return true
	}
	return registryHasBase(registrySet, baseCommand(normalized))
}
