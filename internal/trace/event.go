package trace

// Event hashing constants describe the only canonicalization algorithm this
// package currently accepts for replayable trace evidence.
const (
	// HashAlgSHA256 is the stable digest algorithm for event and payload hashes.
	HashAlgSHA256        = "sha256"
	CanonicalSchemaAlgo  = "json-canonicalization-scheme"
	CanonicalAlgoVersion = "1.0.0"
)

// Canonicalization constants used by Slice A event hashes.
