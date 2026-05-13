package prreview

type Citation struct {
	ContextRefID string `json:"context_ref_id,omitempty"`
	DiffHunkID   string `json:"diff_hunk_id,omitempty"`
	SourceDigest string `json:"source_digest,omitempty"`
	LineStart    int    `json:"line_start,omitempty"`
	LineEnd      int    `json:"line_end,omitempty"`
}

// Ledger is the durable review disposition record synthesized from run output
// and any prior human decisions.
