package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

func runCommandSurface(_ context.Context, args []string, stdout, stderr io.Writer) int {
	// Keep argument rejection in the command runner so the discovery surface
	// remains read-only and cannot silently interpret positional payloads.
	if err := commandSurfaceArgError(args); err != nil {
		fmt.Fprintln(stderr, err)
		return exitUsage
	}
	if err := writeCommandSurfaceJSON(stdout); err != nil {
		return 2
	}
	return 0
}

func commandSurfaceArgError(args []string) error {
	// The command-surface endpoint has no subcommands; any argument is a usage
	// error instead of a partial or filtered discovery request.
	if len(args) == 0 {
		return nil
	}
	return fmt.Errorf("command-surface does not accept arguments: %s", strings.Join(args, " "))
}

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
