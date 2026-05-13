package witness

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

func digestBytes(data []byte) string {
	// Policy key binding uses the digest of the parsed public key bytes, not the
	// original PEM text, so formatting differences do not affect authority.
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

func strongDigest(value string) bool {
	if len(value) < 64 {
		// Customer PKI freshness must bind to at least a SHA-256-sized hex digest.
		return false
	}
	_, err := hex.DecodeString(value)
	return err == nil
}

func defaultString(value, fallback string) string {
	if strings.TrimSpace(value) == "" {
		// Blank profile fields inherit the caller's explicit fallback instead
		// of becoming empty trust-context claims.
		return fallback
	}
	return value
}
