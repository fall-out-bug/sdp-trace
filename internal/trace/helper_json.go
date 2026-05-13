package trace

import (
	"context"
	"encoding/json"
	"os"
)

// ReadJSON decodes any JSON file into dst.
func ReadJSON(ctx context.Context, path string, dst any) error {
	// Context is accepted for call-site symmetry with verifier paths; local file
	// reads remain synchronous.
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	_ = ctx
	return json.Unmarshal(data, dst)
}
