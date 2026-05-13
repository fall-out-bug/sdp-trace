package prreview

import (
	"crypto/sha256"
	"encoding/hex"

	"path/filepath"

	"strings"
)

func copiedInputSafeRef(name, kind, contentType string, data []byte) SafeRef {
	// copiedInputSafeRef keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	digest := sha256.Sum256(data)
	return SafeRef{
		ID:             strings.TrimSuffix(name, filepath.Ext(name)),
		Kind:           kind,
		Ref:            filepath.ToSlash(filepath.Join("inputs", name)),
		DigestSHA256:   hex.EncodeToString(digest[:]),
		ContentType:    contentType,
		RedactionState: RedactionNone,
	}
}
