package demo

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
)

func hashFile(path string) (string, error) {
	// hashFile keeps demo gate evidence explicit and replay-bound.
	// Local rows, contract requirements, witness bindings, protected gates, and artifacts stay separate.
	// This helper renders or aggregates demo evidence; it does not create external proof.
	data, err := os.ReadFile(path)

	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:]), err
}
