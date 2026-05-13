package packet

type BundleEntry struct {
	// Ref is the stable evidence namespace used by packet rows and findings.
	Ref                string   `json:"ref"`
	SourceClass        string   `json:"source_class"`
	Digest             string   `json:"digest,omitempty"`
	RetainedForm       string   `json:"retained_form"`
	RedactionStatus    string   `json:"redaction_status"`
	Resolver           string   `json:"resolver,omitempty"`
	ExpiresAt          string   `json:"expires_at,omitempty"`
	ArtifactAccess     string   `json:"artifact_access,omitempty"`
	ProjectionRole     string   `json:"projection_role,omitempty"`
	EvidenceKind       string   `json:"evidence_kind,omitempty"`
	ObservedComponents []string `json:"observed_components,omitempty"`
	ContradictsRef     string   `json:"contradicts_ref,omitempty"`
	ContradictsRowID   string   `json:"contradicts_row_id,omitempty"`
	Actor              string   `json:"actor,omitempty"`
	WriteAuthority     string   `json:"write_authority,omitempty"`
	GeneratedBy        string   `json:"generated_by,omitempty"`
	SourceCommitState  string   `json:"source_commit_state,omitempty"`
	SourceRef          string   `json:"source_ref,omitempty"`
}
