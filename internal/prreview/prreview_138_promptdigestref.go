package prreview

import (
	"crypto/sha256"
	"encoding/hex"

	"path/filepath"
)

func promptDigestRef(role ReviewRole, data []byte) *SafeRef {
	// promptDigestRef keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	sum := sha256.Sum256(data)
	return &SafeRef{
		ID:             "prompt-" + safeID(role.RoleID),
		Kind:           RefKindPrompt,
		Ref:            "digest-only:" + safeID(filepath.Base(role.PromptTemplateRef)),
		DigestSHA256:   hex.EncodeToString(sum[:]),
		ContentType:    contentType(role.PromptTemplateRef),
		RedactionState: RedactionDigestOnly,
	}
}
