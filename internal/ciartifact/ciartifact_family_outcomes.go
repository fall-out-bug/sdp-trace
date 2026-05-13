package ciartifact

type familyOutcome struct {
	state  string
	code   string
	reason string
	action string
}

var accessResults = map[string]familyOutcome{
	AccessUnsafe:       {StateFail, "unsafe_artifact_output", "artifact family matched a forbidden output-safety class", "Remove unsafe artifact content and regenerate the observation."},
	AccessAbsent:       {StateCannotVerify, "family_absent_in_ci_bundle", "required artifact family is absent from the selected CI bundle", "Upload the required artifact family or mark it outside profile scope."},
	AccessPartial:      {StateCannotVerify, "family_partial_in_ci_bundle", "required artifact family is only partially present", "Upload every required artifact for the selected family."},
	AccessExpired:      {StateCannotVerify, "artifact_expired_before_inspection", "artifact family expired before inspection", "Regenerate CI artifacts or preserve them in an accepted external store."},
	AccessInaccessible: {StateCannotVerify, "artifact_inaccessible", "artifact family could not be accessed under the selected profile", "Provide accessible artifact evidence or mark access not assessed."},
	AccessMalformed:    {StateCannotVerify, "artifact_malformed", "artifact family metadata is malformed", "Fix artifact metadata and rerun observation."},
	AccessCannotVerify: {StateCannotVerify, "artifact_access_cannot_verify", "artifact family access could not be verified", "Provide verifier-readable artifact access metadata."},
}

var bindingResults = map[string]familyOutcome{
	BindingMismatch:     {StateFail, "source_run_binding_mismatch", "artifact family binding contradicts the selected source or run", "Regenerate artifact evidence for the selected source and run."},
	BindingAbsent:       {StateCannotVerify, "source_run_binding_missing", "artifact family lacks selected source or run binding", "Record source and run binding for the selected artifact family."},
	BindingUnverifiable: {StateCannotVerify, "external_artifact_ref_unverifiable", "artifact family binding is unverifiable under the selected profile", "Provide digest-checkable artifact binding evidence."},
}
