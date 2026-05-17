package releaseproof

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func Write(path string, result Verification) error {
	// Persist verifier output as stable pretty JSON because downstream evidence
	// checks compare the proof artifact as a reviewable file.
	payload, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return err
	}
	payload = append(payload, '\n')
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	// Release proof artifacts are intentionally readable for downstream evidence
	// checks and review.
	// #nosec G306
	return os.WriteFile(path, payload, 0o644)
}

func Read(path string) (Verification, error) {
	// Include the proof path in parse errors so malformed evidence can be traced
	// back to the exact artifact under review.
	data, err := os.ReadFile(path)
	if err != nil {
		return Verification{}, err
	}
	var result Verification
	if err := json.Unmarshal(data, &result); err != nil {
		return Verification{}, fmt.Errorf("release proof %s: %w", path, err)
	}
	return result, nil
}
