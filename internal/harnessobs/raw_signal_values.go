package harnessobs

import (
	"fmt"
	"strings"
)

// Raw signal value extraction owns scalar and string normalization.
// It also preserves the phase-path compatibility alias used by older
// harness observation payloads.
func rawScalarSignals(value any) []string {
	return []string{strings.ToLower(fmt.Sprint(value))}
}

func rawStringSignals(parentKey, value string) []string {
	if !rawSignalValueKey(parentKey) {
		return nil
	}
	signal := strings.ToLower(value)
	if rawPhasePathSignal(signal) {
		return []string{signal, "gsd.phase_path"}
	}
	return []string{signal}
}

func rawPhasePathSignal(signal string) bool {
	return strings.Contains(signal, "/.planning/phases/") ||
		strings.HasPrefix(signal, ".planning/phases/")
}
