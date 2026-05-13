package interaction

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
)

func WriteContentBlobs(tracePath string, trace Trace, bodies map[string][]byte) error {
	// WriteContentBlobs keeps interaction trace evidence explicit and source-bound.
	// Relay, import, validation, envelope, summary, and retention states stay separate.
	// This helper records or renders observed interaction data; it does not create external proof.

	for _, event := range trace.Events {
		body, ok := bodies[event.InteractionID]
		if !ok {
			continue
		}
		if err := writeContentBlob(tracePath, event, body); err != nil {
			return err
		}
	}
	return nil
}

func writeContentBlob(tracePath string, event Event, body []byte) error {
	// writeContentBlob keeps interaction trace evidence explicit and source-bound.
	// Relay, import, validation, envelope, summary, and retention states stay separate.
	// This helper records or renders observed interaction data; it does not create external proof.

	sum := sha256.Sum256(body)
	if event.ContentDigest != hex.EncodeToString(sum[:]) {
		return fmt.Errorf("content digest mismatch for %s", event.InteractionID)
	}
	path := contentBlobPath(tracePath, event)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, body, 0o600)
}
