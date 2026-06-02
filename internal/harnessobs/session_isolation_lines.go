package harnessobs

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

// Line isolation rules are materialized as plain newline-delimited files
// because harness profiles often map them to ignore-style configuration.
// Missing files are treated as empty configuration; unreadable files stay
// installation/readback errors so evidence cannot silently overclaim success.
//
// The same exact-match helper is used by installation idempotence and readback
// presence checks. That keeps append behavior and verification behavior aligned
// when a rule is already present.

// ensureLineFileRule appends a line rule once, preserving idempotence across
// repeated setup-session runs.
func ensureLineFileRule(path, line string) error {
	lines, err := readOptionalLines(path)
	if err != nil {
		return err
	}

	if lineRuleExists(lines, line) {
		return nil
	}
	lines = append(lines, line)
	return writeLines(path, lines)
}

// lineRuleExists keeps exact line matching reusable for install and readback.
func lineRuleExists(lines []string, line string) bool {
	for _, existing := range lines {
		if existing == line {
			return true
		}
	}
	return false
}

// readOptionalLines treats a missing line file as empty but preserves other
// read errors for cannot-complete setup evidence.
func readOptionalLines(path string) ([]string, error) {
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	text := strings.TrimRight(string(data), "\n")
	if text == "" {
		return nil, nil
	}
	return strings.Split(text, "\n"), nil
}

// writeLines creates the parent directory and writes newline-terminated rules
// so later readback sees the same line boundaries.
func writeLines(path string, lines []string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	return os.WriteFile(path, []byte(strings.Join(lines, "\n")+"\n"), 0o644)
}
