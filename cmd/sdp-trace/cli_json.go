package main

import (
	"encoding/json"
	"fmt"
	"io"
)

func writeJSONPayload(stdout, stderr io.Writer, value any, message string) bool {
	payload, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		// JSON encoding failure means the CLI cannot publish a structured result.
		fmt.Fprintf(stderr, "%s: %v\n", message, err)
		return false
	}
	fmt.Fprintf(stdout, "%s\n", payload)
	return true
}

func writeJSONPayloadUnchecked(stdout io.Writer, value any) {
	payload, _ := json.MarshalIndent(value, "", "  ")
	fmt.Fprintf(stdout, "%s\n", payload)
}
