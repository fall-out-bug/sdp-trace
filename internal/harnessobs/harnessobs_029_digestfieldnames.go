package harnessobs

var digestFieldNames = map[string]bool{
	"source_digest":     true,
	"validation_digest": true,
	"commit_digest":     true,
	"envelope_digest":   true,
	"payload_digest":    true,
	"sha256":            true,
}
