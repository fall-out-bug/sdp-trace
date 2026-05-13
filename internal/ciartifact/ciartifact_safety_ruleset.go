package ciartifact

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

func defaultSafetyRuleset(input SafetyRuleset) SafetyRuleset {
	// The default safety ruleset is explicit product behavior, not hidden policy.
	// Callers can compare its digest to understand which checks ran.
	if !safeToken(input.ID) {

		input.ID = SafetyRulesetDefault
	}
	if !safeHex(input.SHA256, 64) {

		sum := sha256.Sum256([]byte(defaultSafetyRulesetContent()))
		input.SHA256 = hex.EncodeToString(sum[:])
	}
	return input
}

func defaultSafetyRulesetContent() string {
	// Ruleset content stays deterministic so its digest is stable evidence for the
	// safety rules applied to this observation.

	return strings.Join([]string{
		SafetyRulesetDefault,
		"token_like",
		"jwt_token",
		"private_key",
		"cloud_credential",
		"provider_token",
		"private_artifact_url",
		"private_filesystem_path",
		"prompt_or_model_payload",
		"raw_job_log",
		"high_entropy_secret",
	}, "\n")
}

func defaultString(value, fallback string) string {
	// defaultString keeps CI artifact assessment evidence structured and replayable.
	// It preserves explicit states instead of turning local manifest data into proof.
	// Non-pass outcomes remain tied to a concrete source, family, or safety boundary.
	if strings.TrimSpace(value) == "" {

		return fallback
	}
	return value
}
