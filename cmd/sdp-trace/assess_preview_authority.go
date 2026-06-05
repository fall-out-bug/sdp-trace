package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/authority"
	"io"
)

func runAuthorityAssessPreview(opts *flagSet, stdout io.Writer) int {
	// Authority preview treats the package path as an input pointer only; it does
	// not evaluate authority rules.
	inputs := map[string]string{
		"authority_package": managedInputStatus(opts.stringValue("authority-package")),
	}
	writeJSONPayloadUnchecked(stdout, newAuthorityPreviewReport(inputs))
	return previewInputExitCode(inputs)
}

func newAuthorityPreviewReport(inputs map[string]string) authorityPreviewReport {
	// Authority preview reports package readiness and state vocabulary without
	// evaluating policy effects or accepting raw prompt/model material.
	return authorityPreviewReport{
		Command:         "assess preview",
		SelectedProfile: authority.Profile,
		Inputs:          inputs,
		StateModel:      authorityPreviewStateModel,
		Safety:          authorityPreviewSafety,
		NextActions:     authorityPreviewActions(inputs),
		Claim:           "preview is read-only and does not emit an authority or policy verdict",
	}
}

var authorityPreviewStateModel = map[string]string{
	"authority":   "within_authority,outside_authority,not_assessed,cannot_verify",
	"attribution": "verified,not_assessed,cannot_verify",
	"binding":     "verified,not_assessed,cannot_verify",
}

var authorityPreviewSafety = map[string]string{
	"raw_prompts":       "not_accepted",
	"raw_model_outputs": "not_accepted",
	"credential_refs":   "rejected_as_malformed",
	"policy_effects":    "not_emitted",
}
