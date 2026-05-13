package prreview

func packetVerificationRefs(inputDir string, paths []string) ([]SafeRef, error) {
	// packetVerificationRefs keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.
	refs, err := copyInputs(inputDir, "verification", paths)
	if err != nil {
		return nil, err
	}
	for i := range refs {

		refs[i].Kind = RefKindVerification
	}
	return refs, nil
}
