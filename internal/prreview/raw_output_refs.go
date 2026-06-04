package prreview

import (
	"path/filepath"
)

func digestOnlyRawOutputRef(result ReviewerResult, name, digest string) *SafeRef {
	return rawOutputRef(result, "digest-only:"+name, digest)
}

func retainedRawOutputRef(result ReviewerResult, name, digest string) *SafeRef {
	return rawOutputRef(result, filepath.ToSlash(filepath.Join("raw", name)), digest)
}

func rawOutputRef(result ReviewerResult, ref, digest string) *SafeRef {
	return &SafeRef{ID: "raw-" + safeID(result.ReviewRunID), Kind: RefKindRawOutput, Ref: ref, DigestSHA256: digest, ContentType: ContentText, RedactionState: RedactionDigestOnly}
}
