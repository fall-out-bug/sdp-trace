package main

import (
	"encoding/json"
	"fmt"
	"io"
)

func writeCommandSurfaceJSON(w io.Writer) error {
	surface := buildCommandSurface()
	b, err := json.MarshalIndent(surface, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal command surface: %w", err)
	}
	_, err = w.Write(b)
	if err != nil {
		return fmt.Errorf("write command surface: %w", err)
	}
	_, err = fmt.Fprintln(w)
	return err
}
