package main

import (
	"encoding/json"
	"fmt"
	"io"
)

func writeJSONPayloadUnchecked(stdout io.Writer, value any) {
	payload, _ := json.MarshalIndent(value, "", "  ")
	fmt.Fprintf(stdout, "%s\n", payload)
}
