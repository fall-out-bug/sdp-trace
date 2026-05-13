package prreview

import (
	"fmt"
)

func copiedInputRef(inputDir, prefix string, index int, path string) (SafeRef, error) {
	// copiedInputRef keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	position := index + 1
	name := fmt.Sprintf("%s-%d%s", prefix, position, normalizedExt(path))
	ref, err := copyInput(inputDir, name, path, RefKindDoc, contentType(path))
	if err != nil {
		return SafeRef{}, err
	}
	ref.ID = fmt.Sprintf("%s-%d", prefix, position)
	return ref, nil
}
