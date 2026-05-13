package prreview

import (
	"os"
)

func copyInput(inputDir, name, source, kind, contentType string) (SafeRef, error) {
	// copyInput keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	data, err := os.ReadFile(source)
	if err != nil {
		return SafeRef{}, err
	}
	if err := writeCopiedInput(inputDir, name, data); err != nil {
		return SafeRef{}, err
	}
	return copiedInputSafeRef(name, kind, contentType, data), nil
}
