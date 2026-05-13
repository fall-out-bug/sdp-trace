package query

import (
	"encoding/json"
	"errors"
	"os"
)

func readOptionalPackArtifact(path, role, redactedID string, required bool, target any) (QueryPackInputArtifact, bool, error) {
	// readOptionalPackArtifact keeps query-pack rows source-bound to replayed evidence artifacts.
	// Missing, malformed, redacted, retained, and adapter evidence stay separate.
	// This helper renders derived query rows; it does not create a new verdict.
	if _, err := os.Stat(path); err != nil {
		return optionalArtifactStatResult(err, role, redactedID, required)
	}
	artifact, err := readPackArtifact(path, role, redactedID, required, target)
	return artifact, true, err
}

func optionalArtifactStatResult(err error, role, redactedID string, required bool) (QueryPackInputArtifact, bool, error) {
	// optionalArtifactStatResult keeps query-pack rows source-bound to replayed evidence artifacts.
	// Missing, malformed, redacted, retained, and adapter evidence stay separate.
	// This helper renders derived query rows; it does not create a new verdict.
	if errors.Is(err, os.ErrNotExist) {
		return QueryPackInputArtifact{}, false, nil
	}
	return QueryPackInputArtifact{
		Role:             role,
		PathRedactedID:   redactedID,
		ArtifactRequired: required,
	}, true, err
}

func readPackArtifact(path, role, redactedID string, required bool, target any) (QueryPackInputArtifact, error) {
	// readPackArtifact keeps query-pack rows source-bound to replayed evidence artifacts.
	// Missing, malformed, redacted, retained, and adapter evidence stay separate.
	// This helper renders derived query rows; it does not create a new verdict.
	artifact := QueryPackInputArtifact{
		Role:             role,
		PathRedactedID:   redactedID,
		ArtifactRequired: required,
	}
	payload, err := os.ReadFile(path)
	if err != nil {
		return artifact, err
	}
	artifact.SHA256 = sha256Hex(payload)
	artifact.SchemaVersion = readArtifactSchemaVersion(payload)
	if err := json.Unmarshal(payload, target); err != nil {
		return artifact, err
	}
	return artifact, nil
}
