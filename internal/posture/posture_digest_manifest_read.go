package posture

import (
	"encoding/json"
	"fmt"
	"os"
)

func readDigestManifest(path string) (DigestManifest, error) {
	// readDigestManifest keeps posture export evidence explicit and source-bound.
	// Selection, digest, signal, metric, movement, refusal, and safety states stay separate.
	// This helper aggregates replayed query-pack data; it does not create new proof.
	var manifest DigestManifest
	data, err := os.ReadFile(path)
	if err != nil {
		return manifest, err
	}

	if err := json.Unmarshal(data, &manifest); err != nil {
		return manifest, err
	}
	if manifest.SchemaVersion != DigestManifestSchemaVersion {

		return manifest, fmt.Errorf("unsupported digest manifest schema")
	}

	return manifest, nil
}
