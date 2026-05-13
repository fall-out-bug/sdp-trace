package main

import (
	"slices"
)

func hasGateState(states []string, targets ...string) bool {
	for _, state := range states {
		// Match against the closed state vocabulary selected by the caller.
		if slices.Contains(targets, state) {
			return true
		}
	}
	return false
}
