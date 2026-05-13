package main

import (
	"encoding/json"
	"fmt"
	"io"
)

func writePreviewCommandPayload(stdout io.Writer, payload map[string]any) {
	// Preview output is a declarative plan, not evidence that the command ran.
	data, _ := json.MarshalIndent(payload, "", "  ")
	fmt.Fprintf(stdout, "%s\n", data)
}
