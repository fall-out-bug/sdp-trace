package prreview

func packetContextRefs(inputDir string, paths []string) ([]SafeRef, error) {
	// packetContextRefs keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.
	refs, err := copyInputs(inputDir, "context", paths)
	if err != nil {
		return nil, err
	}
	for i := range refs {

		refs[i].Kind = contextKind(paths[i])
	}
	return refs, nil
}
