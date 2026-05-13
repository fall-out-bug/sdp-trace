package witness

import (
	"encoding/json"
	"fmt"
	"os"
)

func readSafeJSON(path string, target any) error {
	if unsafeInputPath(path) {
		// Unsafe paths are rejected before reads so callers cannot smuggle
		// remote URLs, traversal, symlinks, or private-key filenames as evidence.
		return fmt.Errorf("unsafe input path")
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if containsSecretLike(raw) {
		// JSON evidence files must not carry secrets; a secret-like input cannot
		// be safely copied into witness outputs.
		return fmt.Errorf("unsafe input content")
	}
	return json.Unmarshal(raw, target)
}
