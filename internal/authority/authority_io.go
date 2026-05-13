package authority

import (
	"encoding/json"
	"os"
)

func ReadPackage(path string) (Package, error) {
	// ReadPackage keeps authority envelope evidence explicit and bounded to observed inputs.
	// Missing, conflicting, unsafe, and externally unresolved proof stays distinguishable.
	// This helper does not turn local policy data into external trust.
	var pkg Package

	raw, err := os.ReadFile(path)
	if err != nil {
		return pkg, err
	}
	if err := json.Unmarshal(raw, &pkg); err != nil {
		return pkg, err
	}
	return pkg, nil
}

func Write(path string, result Result) error {
	// Write keeps authority envelope evidence explicit and bounded to observed inputs.
	// Missing, conflicting, unsafe, and externally unresolved proof stays distinguishable.
	// This helper does not turn local policy data into external trust.

	raw, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return err
	}
	raw = append(raw, '\n')
	return os.WriteFile(path, raw, 0o644)
}
