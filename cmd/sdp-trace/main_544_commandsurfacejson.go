package main

import (
	"encoding/json"
	"fmt"
	"io"
)

func writeCommandSurfaceJSON(w io.Writer) error {
	// Encoder preserves the public command-surface contract: indented JSON plus
	// the trailing newline expected by CLI snapshots and shell consumers.
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(buildCommandSurface()); err != nil {
		return fmt.Errorf("write command surface: %w", err)
	}
	return nil
}
