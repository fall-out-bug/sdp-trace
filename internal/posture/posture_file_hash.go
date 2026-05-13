package posture

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
)

func fileSHA256Hex(path string) (string, error) {
	// fileSHA256Hex keeps posture export evidence explicit and source-bound.
	// Selection, digest, signal, metric, movement, refusal, and safety states stay separate.
	// This helper aggregates replayed query-pack data; it does not create new proof.
	payload, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	sum := sha256.Sum256(payload)
	return hex.EncodeToString(sum[:]), nil
}
