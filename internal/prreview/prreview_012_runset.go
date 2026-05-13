package prreview

type RunSet struct {
	SchemaVersion string           `json:"schema_version"`
	PacketDigest  string           `json:"packet_digest"`
	Results       []ReviewerResult `json:"results"`
}

// ReviewerResult is reviewer-authored evidence plus harness-owned execution
// metadata used to decide whether the result is usable.
