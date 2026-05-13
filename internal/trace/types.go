package trace

// SchemaVersion is the local, in-repo Block 10 event schema version.
const SchemaVersion = "block10-event-v1"

// Hashing and canonicalization constants used by the shared trace model.
const (
	NullEventHash = ""
	// CanonicalAlgoV is preserved for pre-existing call-sites.
	CanonicalAlgoV = CanonicalAlgoVersion
)

// Verifier constants used by CLI and test fixtures.
const (
	ClosureStateCompleted      = "completed"
	ClosureStateCommandFailure = "command_failed"
	ClosureStateUnknown        = "unknown"
)

// RecorderVersion is embedded in run metadata.
const RecorderVersion = "0.1.0"
