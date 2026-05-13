package posture

func verifyDigestManifest(manifestPath, queryPackPath string) (string, error) {
	// verifyDigestManifest keeps posture export evidence explicit and source-bound.
	// Selection, digest, signal, metric, movement, refusal, and safety states stay separate.
	// This helper aggregates replayed query-pack data; it does not create new proof.

	expected, err := expectedDigestForQueryPack(manifestPath, queryPackPath)
	if err != nil {
		return "", err
	}

	actual, err := fileSHA256Hex(queryPackPath)
	if err != nil {
		return "", err
	}
	return checkedDigest(actual, expected)
}

func expectedDigestForQueryPack(manifestPath, queryPackPath string) (string, error) {
	// expectedDigestForQueryPack keeps posture export evidence explicit and source-bound.
	// Selection, digest, signal, metric, movement, refusal, and safety states stay separate.
	// This helper aggregates replayed query-pack data; it does not create new proof.

	manifest, err := readDigestManifest(manifestPath)
	if err != nil {
		return "", err
	}
	return digestForQueryPackFromManifest(manifest, queryPackPath)
}

// checkedDigest enforces the evidence boundary where manifest expectation meets artifact reality.
// A digest mismatch crosses from unverifiable into untrusted, distinct from missing or malformed.
func checkedDigest(actual, expected string) (string, error) {
	// checkedDigest keeps posture export evidence explicit and source-bound.
	// Selection, digest, signal, metric, movement, refusal, and safety states stay separate.
	// This helper aggregates replayed query-pack data; it does not create new proof.

	if expected != actual {
		return "", errDigestMismatch
	}
	return actual, nil
}
