package prreview

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
