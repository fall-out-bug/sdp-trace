package query

import (
	"crypto/sha256"
	"encoding/hex"
)

func sha256Hex(payload []byte) string {
	sum := sha256.Sum256(payload)
	return hex.EncodeToString(sum[:])
}
