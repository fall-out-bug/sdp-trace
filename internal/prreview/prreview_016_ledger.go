package prreview

type Ledger struct {
	SchemaVersion string          `json:"schema_version"`
	PacketDigest  string          `json:"packet_digest"`
	Findings      []LedgerFinding `json:"findings"`
}

// LedgerFinding preserves reviewer findings with human disposition state kept
// outside reviewer control.
