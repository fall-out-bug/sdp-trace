package authority

import (
	"regexp"
)

var (
	standardEventTypes = map[string]bool{
		"observe":           true,
		"review":            true,
		"feedback":          true,
		"direct_mutation":   true,
		"commit":            true,
		"merge":             true,
		"ci_run":            true,
		"harness_tool_call": true,
		"gateway_request":   true,
	}
	evidenceRefPattern           = regexp.MustCompile(`^(file:[A-Za-z0-9_./:-]+|artifact:[A-Za-z0-9_.:-]+#[A-Za-z0-9_./#:-]+|git:[A-Fa-f0-9]{40,64}#[A-Za-z0-9_./:-]+|external:[A-Za-z0-9_.:-]+)$`)
	unsafeRefPattern             = regexp.MustCompile(`(?i)(bearer |access_token=|oidc_token|BEGIN [A-Z ]*PRIVATE KEY|raw prompt|raw response|raw_job_log|private_artifact_url)`)
	evidenceRefResolutionReasons = map[string]string{
		"inaccessible": "evidence_ref_inaccessible",
		"malformed":    "evidence_ref_malformed",
		"stale":        "evidence_ref_stale",
	}
)
