package harnessobs

import (
	"strings"
)

// JSON isolation rules update only the permission subtree required by
// setup-session profiles. Existing unrelated settings are preserved, blank or
// missing files become empty objects, and invalid JSON remains an explicit
// installation error.
//
// The implementation deliberately replaces scalar permission/read values with
// objects because the isolation contract needs an addressable pattern map.
// That replacement is local to the required branch; other settings survive the
// write through the shared JSON writer.

// ensureJSONReadDenyRule preserves the existing settings object while forcing a
// single permission.read pattern to deny.
//
//nolint:unparam // path is variable in test coverage and remains explicit for future non-default settings files.
func ensureJSONReadDenyRule(path, pattern string) error {
	config, err := readOptionalJSONObject(path)
	if err != nil {
		return err
	}
	setJSONReadDeny(config, pattern)
	return writeJSON(path, config)
}

// blankJSON identifies empty settings content after whitespace trimming.
func blankJSON(data []byte) bool {
	return strings.TrimSpace(string(data)) == ""
}

// setJSONReadDeny writes only the permission.read branch needed by the
// isolation rule and keeps unrelated top-level settings intact.
func setJSONReadDeny(config map[string]any, pattern string) {
	permission := ensureObject(config, "permission")
	read := ensureObject(permission, "read")

	read[pattern] = "deny"
}

// ensureObject creates or replaces the named child when existing settings use a
// scalar where an object is required for permission rules.
func ensureObject(parent map[string]any, key string) map[string]any {
	if child, ok := parent[key].(map[string]any); ok {
		return child
	}

	child := map[string]any{}
	parent[key] = child
	return child
}
