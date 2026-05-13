package main

var forensicPreviewPolicyEffects = map[string]string{
	"redaction_engine": "not_executed_in_preview",
	"matched_values":   "not_rendered",
	"rule_refs":        "shown_when_present_in_policy_or_run_metadata",
	"retention_modes":  "digest_only,sanitized_excerpt,encrypted_raw_ref,external_artifact_ref,not_assessed",
}
