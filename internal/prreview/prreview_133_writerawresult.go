package prreview

import (
	"crypto/sha256"
	"encoding/hex"
	"os"

	"path/filepath"
)

func writeRawResult(result ReviewerResult, role ReviewRole, rawDir string, output []byte) (ReviewerResult, error) {
	// writeRawResult keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.
	if output != nil {

		name := safeID(result.ReviewRunID) + ".out"
		digest := sha256.Sum256(output)
		if role.RawOutputRetention == RedactionDigestOnly {
			result.RawOutputRef = &SafeRef{ID: "raw-" + safeID(result.ReviewRunID), Kind: RefKindRawOutput, Ref: "digest-only:" + name, DigestSHA256: hex.EncodeToString(digest[:]), ContentType: ContentText, RedactionState: RedactionDigestOnly}
			return result, nil
		}
		path := filepath.Join(rawDir, name)
		if err := os.WriteFile(path, output, 0o600); err != nil {
			return result, err
		}
		result.RawOutputRef = &SafeRef{ID: "raw-" + safeID(result.ReviewRunID), Kind: RefKindRawOutput, Ref: filepath.ToSlash(filepath.Join("raw", name)), DigestSHA256: hex.EncodeToString(digest[:]), ContentType: ContentText, RedactionState: RedactionDigestOnly}
	}
	return result, nil
}
