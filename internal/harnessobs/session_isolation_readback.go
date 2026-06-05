package harnessobs

import "path/filepath"

// Isolation readback is the local verification boundary for setup-session
// mutations. Installation may write a target, but the reported result is only
// pass when the requested line or JSON permission can be read back; absent
// materialization is represented as cannot_verify.
//
// Unsupported rule kinds are errors because no verifier exists for them.
// Existing files that omit the requested rule are not errors: they are local
// evidence that setup could not prove the isolation effect.
//
// Digest attachment is intentionally best-effort after readback. A missing
// digest must not turn an absent rule into pass, and an absent rule must not
// discard the target path, pattern, or reason code that explain the gap.

// verifyIsolationRule converts local file readback into the reported isolation
// result; absent materialization remains cannot_verify, not pass.
func verifyIsolationRule(rule SessionIsolationRule) (SessionIsolationResult, error) {
	result := initialIsolationResult(rule)
	ok, err := isolationRulePresent(rule)
	if err != nil {
		return SessionIsolationResult{}, err
	}
	applyIsolationReadback(&result, ok)
	setIsolationDigest(&result, rule.TargetPath)
	return result, nil
}

// applyIsolationReadback downgrades the result when the requested rule cannot
// be found after installation.
func applyIsolationReadback(result *SessionIsolationResult, ok bool) {
	if !ok {
		result.State = StateCannotVerify
		result.ReasonCode = "isolation_rule_absent"
	}
}

// setIsolationDigest attaches a digest only when the target file exists and can
// be hashed.
func setIsolationDigest(result *SessionIsolationResult, path string) {
	if digest := digestFile(path); digest != "" {
		result.SHA256 = digest
	}
}

// initialIsolationResult builds the optimistic result that readback can still
// downgrade before returning evidence.
func initialIsolationResult(rule SessionIsolationRule) SessionIsolationResult {
	return SessionIsolationResult{
		ID:         rule.ID,
		Kind:       rule.Kind,
		TargetPath: filepath.ToSlash(rule.TargetPath),
		Pattern:    rule.Pattern,
		State:      StatePass,
		ReasonCode: "isolation_rule_verified",
	}
}
