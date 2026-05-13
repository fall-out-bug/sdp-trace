package export

import (
	"encoding/json"
	"os"
)

// Write writes a deterministic JSON bundle file.
func WriteBundle(path string, bundle AuditBundle) error {
	return writeBundle(path, bundle)
}

// Read reads a prebuilt bundle.
func Read(path string) (AuditBundle, error) {
	return readBundle(path)
}

func writeBundle(path string, bundle AuditBundle) error {
	// Indented JSON keeps exported evidence readable in review diffs.
	data, err := json.MarshalIndent(bundle, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

func readBundle(path string) (AuditBundle, error) {
	// Reading a bundle is shape-only; trust still comes from replaying the
	// referenced run and verifier artifacts.
	data, err := os.ReadFile(path)
	if err != nil {
		return AuditBundle{}, err
	}
	var bundle AuditBundle
	if err := json.Unmarshal(data, &bundle); err != nil {
		return AuditBundle{}, err
	}
	return bundle, nil
}
