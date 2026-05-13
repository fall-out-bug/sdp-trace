package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestProofSummaryShapeExampleCannotClaimPass(t *testing.T) {
	payload, err := os.ReadFile(filepath.Join("..", "..", "schema", "proof-summary.schema.json"))
	if err != nil {
		t.Fatalf("read proof-summary schema: %v", err)
	}
	var schema map[string]any
	if err := json.Unmarshal(payload, &schema); err != nil {
		t.Fatalf("decode proof-summary schema: %v", err)
	}
	branch := proofSummaryUntrustedShapeBranch(t, schema)
	thenProps := schemaMapAt(t, schemaMapAt(t, branch, "then"), "properties")
	result := schemaMapAt(t, thenProps, "result")
	enum := stringSliceAt(t, result, "enum")
	if !sameStrings(enum, []string{"cannot_verify", "not_assessed"}) {
		t.Fatalf("untrusted shape result enum = %v", enum)
	}
	notAssessed := schemaMapAt(t, thenProps, "not_assessed")
	if notAssessed["minItems"] != float64(1) {
		t.Fatalf("untrusted shape not_assessed minItems = %v", notAssessed["minItems"])
	}
}

func proofSummaryUntrustedShapeBranch(t *testing.T, schema map[string]any) map[string]any {
	t.Helper()
	for _, item := range schemaArrayAt(t, schema, "allOf") {
		branch, ok := item.(map[string]any)
		if !ok {
			continue
		}
		if props, ok := schemaMapAt(t, schemaMapAt(t, branch, "if"), "properties")["artifact_role"].(map[string]any); ok && props["const"] == "untrusted_shape_example" {
			return branch
		}
	}
	t.Fatalf("untrusted shape branch not found")
	return nil
}

func schemaMapAt(t *testing.T, source map[string]any, key string) map[string]any {
	t.Helper()
	value, ok := source[key].(map[string]any)
	if !ok {
		t.Fatalf("%s is %T, want object", key, source[key])
	}
	return value
}

func schemaArrayAt(t *testing.T, source map[string]any, key string) []any {
	t.Helper()
	value, ok := source[key].([]any)
	if !ok {
		t.Fatalf("%s is %T, want array", key, source[key])
	}
	return value
}

func stringSliceAt(t *testing.T, source map[string]any, key string) []string {
	t.Helper()
	values := schemaArrayAt(t, source, key)
	out := make([]string, 0, len(values))
	for _, value := range values {
		text, ok := value.(string)
		if !ok {
			t.Fatalf("%s item is %T, want string", key, value)
		}
		out = append(out, text)
	}
	return out
}

func sameStrings(left, right []string) bool {
	if len(left) != len(right) {
		return false
	}
	for i := range left {
		if left[i] != right[i] {
			return false
		}
	}
	return true
}
