package query

import "sort"

type packBuilder struct {
	inputs   packInputs
	rows     map[string][]QueryPackRow
	counters map[string]int
}

func newPackBuilder(inputs packInputs) *packBuilder {
	// newPackBuilder keeps query-pack rows source-bound to replayed evidence artifacts.
	// Missing, malformed, redacted, retained, and adapter evidence stay separate.
	// This helper renders derived query rows; it does not create a new verdict.
	rows := map[string][]QueryPackRow{}
	for _, queryName := range queryOrder {
		rows[queryName] = []QueryPackRow{}
	}
	return &packBuilder{inputs: inputs, rows: rows, counters: map[string]int{}}
}

func (b *packBuilder) result() QueryPackResult {
	// result keeps query-pack rows source-bound to replayed evidence artifacts.
	// Missing, malformed, redacted, retained, and adapter evidence stay separate.
	// This helper renders derived query rows; it does not create a new verdict.
	return QueryPackResult{
		SchemaVersion:    QueryPackSchemaVersion,
		QueryPackID:      QueryPackForensicsBasic,
		QueryPackVersion: "v1",
		RunID:            b.inputs.run.RunID,
		RunNonce:         b.inputs.run.RunNonce,
		SourceBaseline:   b.inputs.run.SourceBaseline,
		InputArtifacts:   b.inputArtifacts(),
		QueryRows:        b.rows,
		OutputSafety:     b.outputSafety(),
	}
}

func (b *packBuilder) inputArtifacts() []QueryPackInputArtifact {
	// inputArtifacts keeps query-pack rows source-bound to replayed evidence artifacts.
	// Missing, malformed, redacted, retained, and adapter evidence stay separate.
	// This helper renders derived query rows; it does not create a new verdict.
	artifacts := []QueryPackInputArtifact{b.inputs.runArtifact}
	artifacts = appendObservedArtifact(artifacts, b.inputs.forensicArtifact)
	artifacts = appendObservedArtifact(artifacts, b.inputs.adapterArtifact)
	sort.Slice(artifacts, func(i, j int) bool { return artifacts[i].Role < artifacts[j].Role })
	return artifacts
}

func appendObservedArtifact(artifacts []QueryPackInputArtifact, artifact *QueryPackInputArtifact) []QueryPackInputArtifact {
	// appendObservedArtifact keeps query-pack rows source-bound to replayed evidence artifacts.
	// Missing, malformed, redacted, retained, and adapter evidence stay separate.
	// This helper renders derived query rows; it does not create a new verdict.
	if artifact == nil {
		return artifacts
	}
	return append(artifacts, *artifact)
}

func (b *packBuilder) outputSafety() QueryPackOutputSafety {
	// outputSafety keeps query-pack rows source-bound to replayed evidence artifacts.
	// Missing, malformed, redacted, retained, and adapter evidence stay separate.
	// This helper renders derived query rows; it does not create a new verdict.
	return QueryPackOutputSafety{
		RedactionPolicyDigest:          b.inputs.run.RedactionDigest,
		VerifiedAbsentSensitiveClasses: sensitiveClasses(),
	}
}
