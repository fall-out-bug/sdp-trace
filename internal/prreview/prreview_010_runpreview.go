package prreview

type RunPreview struct {
	SchemaVersion string        `json:"schema_version"`
	PacketDigest  string        `json:"packet_digest"`
	Roles         []PreviewRole `json:"roles"`
}

// PreviewRole is the per-role portion of a preview, limited to digests and
// runner metadata that are safe to show before execution.
