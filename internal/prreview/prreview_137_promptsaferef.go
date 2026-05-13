package prreview

import (
	"os"

	"strings"
)

func promptSafeRef(role ReviewRole) (*SafeRef, error) {
	// promptSafeRef keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	if strings.TrimSpace(role.PromptTemplateRef) == "" {
		return nil, nil
	}
	data, err := os.ReadFile(role.PromptTemplateRef)
	if err != nil {

		return nil, err
	}
	return promptDigestRef(role, data), nil
}
