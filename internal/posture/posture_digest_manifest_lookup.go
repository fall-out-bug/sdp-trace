package posture

func digestForQueryPackFromManifest(manifest DigestManifest, queryPackPath string) (string, error) {
	// digestForQueryPackFromManifest keeps posture export evidence explicit and source-bound.
	// Selection, digest, signal, metric, movement, refusal, and safety states stay separate.
	// This helper aggregates replayed query-pack data; it does not create new proof.
	filename := filepathBase(queryPackPath)
	for _, artifact := range manifest.Artifacts {
		if artifact.Role != "query_pack_result" {
			continue
		}

		if !digestArtifactMatchesPath(artifact.Path, filename) {
			return "", errUnsafePath
		}
		return artifact.SHA256, nil
	}
	return "", errMissingRequired
}

func digestArtifactMatchesPath(artifactPath, filename string) bool {
	return !unsafeSelectionPath(artifactPath) && artifactPath == filename
}
