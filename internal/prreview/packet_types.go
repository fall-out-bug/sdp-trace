package prreview

// Packet is the frozen review input boundary: reviewers assess the copied refs,
// not the mutable files those refs came from.
type Packet struct {
	SchemaVersion     string             `json:"schema_version"`
	PacketID          string             `json:"packet_id"`
	PacketDigest      string             `json:"packet_digest"`
	RepoID            string             `json:"repo_id"`
	ChangeRef         string             `json:"change_ref"`
	BaseCommit        string             `json:"base_commit"`
	HeadCommit        string             `json:"head_commit"`
	DiffRef           SafeRef            `json:"diff_ref"`
	MetadataRef       *SafeRef           `json:"metadata_ref,omitempty"`
	ContextRefs       []SafeRef          `json:"context_refs"`
	VerificationRefs  []SafeRef          `json:"verification_refs"`
	CIState           string             `json:"ci_state"`
	CreatedAt         string             `json:"created_at"`
	CreatedBy         string             `json:"created_by"`
	RedactionState    string             `json:"redaction_state"`
	UnavailableFields []UnavailableField `json:"unavailable_fields,omitempty"`
}

// SafeRef records enough provenance to replay a packet input without treating
// the referenced content as an approval or proof claim by itself.
type SafeRef struct {
	ID             string `json:"id"`
	Kind           string `json:"kind"`
	Ref            string `json:"ref"`
	DigestSHA256   string `json:"digest_sha256"`
	ContentType    string `json:"content_type"`
	RedactionState string `json:"redaction_state"`
}

// UnavailableField makes missing optional packet inputs explicit so absence
// cannot be mistaken for passing evidence.
type UnavailableField struct {
	Field  string `json:"field"`
	State  string `json:"state"`
	Reason string `json:"reason"`
}
