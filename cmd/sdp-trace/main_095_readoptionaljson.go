package main

import (
	"encoding/json"
	"os"
	"strings"
)

func readOptionalJSON(path string, target any) error {
	if strings.TrimSpace(path) == "" {
		// Empty optional path means "no supplemental evidence".
		return nil
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(raw, target)
}
