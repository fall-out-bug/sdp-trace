package query

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestQueryPackSchemaMirrorsImplementationEnums(t *testing.T) {
	schemaPath := filepath.Join("..", "..", "schema", "forensics-query-pack-result.schema.json")
	payload, err := os.ReadFile(schemaPath)
	if err != nil {
		t.Fatalf("read query-pack schema: %v", err)
	}
	var schema map[string]any
	if err := json.Unmarshal(payload, &schema); err != nil {
		t.Fatalf("decode query-pack schema: %v", err)
	}

	properties := objectAt(t, schema, "properties")
	assertConst(t, properties, "schema_version", QueryPackSchemaVersion)
	assertConst(t, properties, "query_pack_id", QueryPackForensicsBasic)

	queryRow := objectAt(t, objectAt(t, schema, "$defs"), "queryRow")
	rowProperties := objectAt(t, queryRow, "properties")
	assertEnum(t, rowProperties, "query", queryOrder)
	assertEnum(t, rowProperties, "evidence_state", []string{
		RowStatePresent,
		RowStateIssueObserved,
		RowStateNotAssessed,
		RowStateCannotVerify,
		RowStateMissingTelemetry,
		RowStateUnsupported,
		RowStateNotIntegrated,
		RowStateRetentionLimited,
	})
	assertEnum(t, rowProperties, "evidence_family", []string{
		EvidenceFamilyRunChain,
		EvidenceFamilyWitness,
		EvidenceFamilyRetention,
		EvidenceFamilyRedaction,
		EvidenceFamilyAdapterCapture,
		EvidenceFamilyTask,
		EvidenceFamilyCommand,
		EvidenceFamilyFileMutations,
		EvidenceFamilyTest,
		EvidenceFamilySupersession,
		EvidenceFamilyClaim,
		EvidenceFamilyInputArtifact,
	})
}

func objectAt(t *testing.T, source map[string]any, key string) map[string]any {
	t.Helper()
	value, ok := source[key].(map[string]any)
	if !ok {
		t.Fatalf("%s is not an object in %+v", key, source)
	}
	return value
}

func assertConst(t *testing.T, properties map[string]any, key string, expected string) {
	t.Helper()
	property := objectAt(t, properties, key)
	actual, ok := property["const"].(string)
	if !ok {
		t.Fatalf("%s.const is not a string: %+v", key, property)
	}
	if actual != expected {
		t.Fatalf("%s.const = %q expected %q", key, actual, expected)
	}
}

func assertEnum(t *testing.T, properties map[string]any, key string, expected []string) {
	t.Helper()
	property := objectAt(t, properties, key)
	rawEnum, ok := property["enum"].([]any)
	if !ok {
		t.Fatalf("%s.enum is not an array: %+v", key, property)
	}
	actual := make([]string, 0, len(rawEnum))
	for _, value := range rawEnum {
		asString, ok := value.(string)
		if !ok {
			t.Fatalf("%s.enum contains non-string value: %+v", key, rawEnum)
		}
		actual = append(actual, asString)
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("%s.enum = %+v expected %+v", key, actual, expected)
	}
}
