package harnessobs

import "errors"

type isolationPresenceChecker func(path, pattern string) (bool, error)

var isolationPresenceCheckers = map[string]isolationPresenceChecker{
	"ignore_line":    lineIsolationRulePresent,
	"json_read_deny": jsonReadDenyRulePresent,
}

// Isolation presence checks are read-only probes used after installation.
// They return false for absent rules, preserve read errors, and reject unknown
// rule kinds because no evidence contract exists for them.

// isolationRulePresent routes readback by rule kind and rejects unsupported
// kinds instead of inventing a verification state.
func isolationRulePresent(rule SessionIsolationRule) (bool, error) {
	checker, ok := isolationPresenceCheckers[rule.Kind]
	if !ok {
		return false, errors.New("unsupported isolation rule kind")
	}
	return checker(rule.TargetPath, rule.Pattern)
}

// lineIsolationRulePresent checks exact line equality so partial pattern
// matches cannot satisfy ignore-line evidence.
func lineIsolationRulePresent(path, pattern string) (bool, error) {
	lines, err := readOptionalLines(path)
	if err != nil {
		return false, err
	}

	return lineRuleExists(lines, pattern), nil
}

// jsonReadDenyRulePresent reads only the permission subtree needed for the
// read-deny proof while sharing the existing JSON reader helpers.
func jsonReadDenyRulePresent(path, pattern string) (bool, error) {
	// Decode only the subtree needed for this isolation proof.
	var config struct {
		Permission struct {
			Read map[string]string `json:"read"`
		} `json:"permission"`
	}
	if err := readExistingJSON(path, &config); err != nil {
		return false, err
	}
	return config.Permission.Read[pattern] == "deny", nil
}
