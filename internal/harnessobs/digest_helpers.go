package harnessobs

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
)

// sha256Hex is the shared representation for local digests; it is not an
// authority claim unless a caller binds it to verified source evidence.
func sha256Hex(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

// digestLine ignores the mutable source_digest field when raw events can be
// parsed as JSON, then falls back to digesting the original bytes.
func digestLine(line []byte) string {
	if canonical, ok := canonicalSourceDigestLine(line); ok {
		return sha256Hex(canonical)
	}
	return sha256Hex(line)
}

// canonicalSourceDigestLine keeps source digest recomputation stable by
// zeroing the previous digest before canonical JSON encoding.
func canonicalSourceDigestLine(line []byte) ([]byte, bool) {
	var raw map[string]any
	if err := json.Unmarshal(line, &raw); err != nil {
		return nil, false
	}
	raw["source_digest"] = ""

	canonical, err := json.Marshal(raw)
	if err != nil {
		return nil, false
	}
	return canonical, true
}

// digestFile deliberately returns an empty digest on read failure so callers can
// keep degraded source state explicit instead of overclaiming file evidence.
func digestFile(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return sha256Hex(data)
}

// validationDigest excludes the digest field itself so repeated validation
// rendering is stable.
func validationDigest(validation Validation) string {
	copy := validation
	copy.ValidationDigest = ""

	data, _ := json.Marshal(copy)
	return sha256Hex(data)
}

// digestCommand records the exact command vector retained by the setup or
// collection flow, not a shell-normalized command string.
func digestCommand(command []string) string {
	data, _ := json.Marshal(command)
	return sha256Hex(data)
}
