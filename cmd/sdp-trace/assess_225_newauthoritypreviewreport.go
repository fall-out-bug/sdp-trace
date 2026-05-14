package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/authority"
)

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
