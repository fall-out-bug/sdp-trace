package main

import (
	"io"

	"github.com/fall_out_bug/sdp-trace/internal/adaptercapture"
	"github.com/fall_out_bug/sdp-trace/internal/authority"
	"github.com/fall_out_bug/sdp-trace/internal/ciartifact"
	"github.com/fall_out_bug/sdp-trace/internal/forensic"
	"github.com/fall_out_bug/sdp-trace/internal/managed"
)

type assessmentExplainHandler func(string, io.Writer, io.Writer) int

var assessmentExplainHandlers = map[string]assessmentExplainHandler{
	// Schema versions select typed explainers so spoofed profile names cannot
	// redirect artifact interpretation.
	adaptercapture.SchemaVersion:  explainTypedAssessment[adaptercapture.AssessmentResult](explainAdapterCaptureAssessment),
	managed.SchemaVersion:         explainTypedAssessment[managed.AssessmentResult](explainManagedAssessment),
	forensic.SchemaVersion:        explainTypedAssessment[forensic.AssessmentResult](explainForensicAssessment),
	ciartifact.SchemaVersion:      explainTypedAssessment[ciartifact.ObservationResult](explainCIArtifactObservation),
	authority.ResultSchemaVersion: explainTypedAssessment[authority.Result](explainAuthorityEvaluation),
}
