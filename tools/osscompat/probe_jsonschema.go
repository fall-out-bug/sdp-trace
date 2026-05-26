package main

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"
)

// runCheckJSONSchema runs check-jsonschema against a JSON file.
func runCheckJSONSchema(ctx context.Context, schemaPath, jsonPath string) ([]byte, error) {
	out, err := runExternalTool(ctx, "check-jsonschema", "--schemafile", schemaPath, jsonPath)
	return out, err
}

// runJSONSchemaFixtures validates a checked fixture against the schema.
func runJSONSchemaFixtures() (verifierState, string) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	root := repoRoot()
	// Fixture validation checks the richer profile schema, not the live wrap
	// manifest schema.
	out, err := runCheckJSONSchema(ctx,
		filepath.Join(root, "schema/flight-recorder-run.schema.json"),
		filepath.Join(root, "examples/flight-recorder/local-positive/run.json"),
	)
	if err != nil {
		return stateFail, fmt.Sprintf("fixture validation failed: %v\n%s", err, strings.TrimSpace(string(out)))
	}
	return statePass, "fixture validates against schema"
}
