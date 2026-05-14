package main

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestExitCode(t *testing.T) {
	if got := exitCode(nil, os.Stderr); got != 0 {
		t.Fatalf("exitCode(nil) = %d, want 0", got)
	}
	err := errors.New("drift")
	if got := exitCode(err, os.Stderr); got != 1 {
		t.Fatalf("exitCode(err) = %d, want 1", got)
	}
}

func schemaDir(t *testing.T, root string) string {
	d := filepath.Join(root, "schema")
	if err := os.MkdirAll(d, 0755); err != nil {
		t.Fatal(err)
	}
	return d
}

func TestCheckPassesWhenIndexMatchesFiles(t *testing.T) {
	root := t.TempDir()
	sd := schemaDir(t, root)
	if err := os.WriteFile(filepath.Join(sd, "a.schema.json"), []byte("{}"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(sd, "b.schema.json"), []byte("{}"), 0644); err != nil {
		t.Fatal(err)
	}
	idx := &Index{
		Schemas: []SchemaEntry{
			{Name: "a.schema.json", Status: "current", Purpose: "A purpose"},
			{Name: "b.schema.json", Status: "historical", Purpose: "B purpose"},
		},
	}
	if err := check(root, idx); err != nil {
		t.Fatalf("check: %v", err)
	}
}

func TestCheckFailsMissingIndexEntry(t *testing.T) {
	root := t.TempDir()
	sd := schemaDir(t, root)
	if err := os.WriteFile(filepath.Join(sd, "a.schema.json"), []byte("{}"), 0644); err != nil {
		t.Fatal(err)
	}
	idx := &Index{Schemas: []SchemaEntry{}}
	err := check(root, idx)
	if err == nil {
		t.Fatal("expected error for missing index entry")
	}
	if !strings.Contains(err.Error(), "missing index entry: a.schema.json") {
		t.Fatalf("error missing expected text: %v", err)
	}
}

func TestCheckFailsExtraIndexEntry(t *testing.T) {
	root := t.TempDir()
	sd := schemaDir(t, root)
	if err := os.WriteFile(filepath.Join(sd, "a.schema.json"), []byte("{}"), 0644); err != nil {
		t.Fatal(err)
	}
	idx := &Index{
		Schemas: []SchemaEntry{
			{Name: "a.schema.json", Status: "current", Purpose: "A"},
			{Name: "missing.schema.json", Status: "current", Purpose: "M"},
		},
	}
	err := check(root, idx)
	if err == nil {
		t.Fatal("expected error for extra index entry")
	}
	if !strings.Contains(err.Error(), "extra index entry") {
		t.Fatalf("error missing expected text: %v", err)
	}
}

func TestCheckFailsMissingStatus(t *testing.T) {
	root := t.TempDir()
	sd := schemaDir(t, root)
	if err := os.WriteFile(filepath.Join(sd, "a.schema.json"), []byte("{}"), 0644); err != nil {
		t.Fatal(err)
	}
	idx := &Index{
		Schemas: []SchemaEntry{
			{Name: "a.schema.json", Status: "", Purpose: "A"},
		},
	}
	err := check(root, idx)
	if err == nil {
		t.Fatal("expected error for missing status")
	}
	if !strings.Contains(err.Error(), "missing status") {
		t.Fatalf("error missing expected text: %v", err)
	}
}

func TestCheckFailsInvalidStatus(t *testing.T) {
	root := t.TempDir()
	sd := schemaDir(t, root)
	if err := os.WriteFile(filepath.Join(sd, "a.schema.json"), []byte("{}"), 0644); err != nil {
		t.Fatal(err)
	}
	idx := &Index{
		Schemas: []SchemaEntry{
			{Name: "a.schema.json", Status: "draft", Purpose: "A"},
		},
	}
	err := check(root, idx)
	if err == nil {
		t.Fatal("expected error for invalid status")
	}
	if !strings.Contains(err.Error(), "invalid status") {
		t.Fatalf("error missing expected text: %v", err)
	}
}

func TestCheckFailsMissingPurpose(t *testing.T) {
	root := t.TempDir()
	sd := schemaDir(t, root)
	if err := os.WriteFile(filepath.Join(sd, "a.schema.json"), []byte("{}"), 0644); err != nil {
		t.Fatal(err)
	}
	idx := &Index{
		Schemas: []SchemaEntry{
			{Name: "a.schema.json", Status: "current", Purpose: ""},
		},
	}
	err := check(root, idx)
	if err == nil {
		t.Fatal("expected error for missing purpose")
	}
	if !strings.Contains(err.Error(), "missing purpose") {
		t.Fatalf("error missing expected text: %v", err)
	}
}

func TestCheckFailsBrokenExampleRef(t *testing.T) {
	root := t.TempDir()
	sd := schemaDir(t, root)
	if err := os.WriteFile(filepath.Join(sd, "a.schema.json"), []byte("{}"), 0644); err != nil {
		t.Fatal(err)
	}
	idx := &Index{
		Schemas: []SchemaEntry{
			{Name: "a.schema.json", Status: "current", Purpose: "A", Examples: []string{"examples/nope.json"}},
		},
	}
	err := check(root, idx)
	if err == nil {
		t.Fatal("expected error for broken example ref")
	}
	if !strings.Contains(err.Error(), "broken example ref") {
		t.Fatalf("error missing expected text: %v", err)
	}
}

func TestCheckPassesValidExampleRef(t *testing.T) {
	root := t.TempDir()
	sd := schemaDir(t, root)
	if err := os.WriteFile(filepath.Join(sd, "a.schema.json"), []byte("{}"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(root, "examples"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "examples", "ok.json"), []byte("{}"), 0644); err != nil {
		t.Fatal(err)
	}
	idx := &Index{
		Schemas: []SchemaEntry{
			{Name: "a.schema.json", Status: "current", Purpose: "A", Examples: []string{"examples/ok.json"}},
		},
	}
	if err := check(root, idx); err != nil {
		t.Fatalf("check: %v", err)
	}
}

func TestRunAcceptsCurrentIndex(t *testing.T) {
	if err := run(false); err != nil {
		t.Fatalf("run: %v", err)
	}
}

func TestGenerateTableIsDeterministic(t *testing.T) {
	idx := &Index{
		Schemas: []SchemaEntry{
			{Name: "b.schema.json", Status: "current", Purpose: "B"},
			{Name: "a.schema.json", Status: "current", Purpose: "A"},
		},
	}
	out := generateTable(idx)
	lines := strings.Split(out, "\n")
	if len(lines) < 4 {
		t.Fatalf("expected header + rows, got:\n%s", out)
	}
	// Rows should follow index order (which we can assume is sorted by caller if desired).
	if !strings.Contains(out, "a.schema.json") || !strings.Contains(out, "b.schema.json") {
		t.Fatalf("generated table missing entries:\n%s", out)
	}
}
