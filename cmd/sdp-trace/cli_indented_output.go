package main

import (
	"encoding/json"
	"fmt"
	"io"
)

func writeIndentedPayload(stdout io.Writer, value any) {
	payload, _ := json.MarshalIndent(value, "", "  ")
	// Callers only pass values that have already been constructed for output.
	fmt.Fprintf(stdout, "%s\n", payload)
}
