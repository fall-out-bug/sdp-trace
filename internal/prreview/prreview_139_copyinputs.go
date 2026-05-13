package prreview

func copyInputs(inputDir, prefix string, paths []string) ([]SafeRef, error) {
	// copyInputs keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	refs := make([]SafeRef, 0, len(paths))
	for i, path := range paths {
		ref, err := copiedInputRef(inputDir, prefix, i, path)
		if err != nil {
			return nil, err
		}
		refs = append(refs, ref)
	}
	return refs, nil
}
