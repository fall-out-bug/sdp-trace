package main

var ciArtifactPreviewStateModel = map[string]string{
	"top_level": "pass,fail,cannot_verify,not_assessed",
	"producer":  "ci_uploaded,checked_in,local_generated,agent_reported,harness_observed,external_artifact_ref,not_assessed",
	"access":    "present,absent,partial,expired,inaccessible,malformed,unsafe,not_assessed,cannot_verify",
}
