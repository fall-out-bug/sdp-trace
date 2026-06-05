package prreview

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"strings"
)

// Prompt refs are digest-only evidence.
//
// Review packets must not embed prompt template contents, but preview and run
// outputs still need a stable digest and content type. Missing prompt refs stay
// empty; unreadable prompt refs are surfaced to the run path as cannot_verify.
func promptDigestForRole(role ReviewRole) string {
	promptRef, _ := promptSafeRef(role)
	if promptRef == nil {
		return ""
	}
	return promptRef.DigestSHA256
}

func promptSafeRef(role ReviewRole) (*SafeRef, error) {
	if strings.TrimSpace(role.PromptTemplateRef) == "" {
		return nil, nil
	}
	data, err := os.ReadFile(role.PromptTemplateRef)
	if err != nil {
		return nil, err
	}
	return promptDigestRef(role, data), nil
}

func promptDigestRef(role ReviewRole, data []byte) *SafeRef {
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
