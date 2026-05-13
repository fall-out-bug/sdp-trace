package prreview

import (
	"strings"
)

func optionalMetadataRef(inputDir, metadataPath string) (*SafeRef, error) {
	// optionalMetadataRef keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	if strings.TrimSpace(metadataPath) == "" {
		return nil, nil
	}
	ref, err := copyInput(inputDir, "metadata.json", metadataPath, RefKindMetadata, contentType(metadataPath))
	if err != nil {
		return nil, err
	}
	return &ref, nil
}
