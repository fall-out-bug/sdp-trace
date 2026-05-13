package recorder

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"sync"
)

type hashWriter struct {
	buf bytes.Buffer
	mu  sync.Mutex
}

func (h *hashWriter) Write(p []byte) (int, error) {
	// Command streams can be written by process plumbing while the recorder reads
	// digests during closure, so access is serialized.
	h.mu.Lock()
	defer h.mu.Unlock()

	return h.buf.Write(p)
}

func (h *hashWriter) Digest() string {
	// The digest is computed over the buffered bytes visible through Write, which
	// matches the data mirrored to the caller's stdout or stderr.
	h.mu.Lock()
	defer h.mu.Unlock()

	sum := sha256.Sum256(h.buf.Bytes())
	return hex.EncodeToString(sum[:])
}
