package prreview

import "strings"

func copyDiffRef(inputDir, diffPath string) (SafeRef, error) {
	return copyInput(inputDir, "diff.patch", diffPath, RefKindDiff, ContentUnifiedDiff)
}

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
