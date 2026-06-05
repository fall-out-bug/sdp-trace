package prreview

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func Copy(r io.Reader, w io.Writer) error {
	_, err := io.Copy(w, r)
	return err
}

func readPacketRef(packetDir string, ref SafeRef) ([]byte, error) {
	// Read and digest verification are intentionally adjacent.
	// The prompt must not include bytes that failed path safety.
	// The prompt must not include bytes whose digest changed after packet build.
	// Both failures collapse to prompt_evidence_cannot_verify for callers.
	path, err := packetRefPath(packetDir, ref.Ref)
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", errPromptEvidenceCannotVerify, ref.ID)
	}
	if !digestMatches(data, ref.DigestSHA256) {
		return nil, fmt.Errorf("%w: %s", errPromptEvidenceCannotVerify, ref.ID)
	}
	return data, nil
}

func packetRefPath(packetDir, ref string) (string, error) {
	// Packet refs are repository-neutral relative paths. Absolute paths and
	// parent traversal would make CI prompts depend on runner filesystem state.
	cleanRef := filepath.Clean(filepath.FromSlash(ref))
	if filepath.IsAbs(cleanRef) || strings.HasPrefix(cleanRef, ".."+string(filepath.Separator)) || cleanRef == ".." {
		return "", fmt.Errorf("%w: unsafe_ref", errPromptEvidenceCannotVerify)
	}
	return filepath.Join(packetDir, cleanRef), nil
}

func digestMatches(data []byte, digest string) bool {
	// SafeRef digests are stored without the sha256: prefix in packet refs.
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:]) == digest
}
