package prreview

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
)

// writeRawResult records a digest-bound raw output reference and persists bytes
// only when the role does not request digest-only retention.
func writeRawResult(result ReviewerResult, role ReviewRole, rawDir string, output []byte) (ReviewerResult, error) {
	if output == nil {
		return result, nil
	}
	name := rawOutputName(result)
	digest := rawOutputDigest(output)
	if role.RawOutputRetention == RedactionDigestOnly {
		result.RawOutputRef = digestOnlyRawOutputRef(result, name, digest)
		return result, nil
	}
	if err := writeRawOutputFile(rawDir, name, output); err != nil {
		return result, err
	}
	result.RawOutputRef = retainedRawOutputRef(result, name, digest)
	return result, nil
}

// rawOutputName keeps raw output filenames derived from safe run IDs.
func rawOutputName(result ReviewerResult) string {
	return safeID(result.ReviewRunID) + ".out"
}

// rawOutputDigest returns the hex SHA-256 digest recorded in raw output refs.
func rawOutputDigest(output []byte) string {
	digest := sha256.Sum256(output)
	return hex.EncodeToString(digest[:])
}

// writeRawOutputFile persists retained raw output with private file mode.
func writeRawOutputFile(rawDir, name string, output []byte) error {
	return os.WriteFile(filepath.Join(rawDir, name), output, 0o600)
}
