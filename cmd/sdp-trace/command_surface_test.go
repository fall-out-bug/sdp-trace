package main

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestCommandSurfaceJSONIsDeterministic(t *testing.T) {
	var out bytes.Buffer
	if err := writeCommandSurfaceJSON(&out); err != nil {
		t.Fatalf("writeCommandSurfaceJSON: %v", err)
	}
	var first, second commandSurfaceSchema
	if err := json.Unmarshal(out.Bytes(), &first); err != nil {
		t.Fatalf("unmarshal first: %v", err)
	}
	out.Reset()
	if err := writeCommandSurfaceJSON(&out); err != nil {
		t.Fatalf("writeCommandSurfaceJSON second: %v", err)
	}
	if err := json.Unmarshal(out.Bytes(), &second); err != nil {
		t.Fatalf("unmarshal second: %v", err)
	}
	if !jsonEqual(first, second) {
		t.Fatal("command surface output is not deterministic")
	}
}

func TestCommandSurfaceCoversAllHandlers(t *testing.T) {
	surface := buildCommandSurface()
	if got, want := len(surface.Commands), len(commandHandlers); got != want {
		t.Fatalf("registry has %d commands, handlers has %d", got, want)
	}
}

func TestCommandSurfaceHasSchemaVersion(t *testing.T) {
	var out bytes.Buffer
	if err := writeCommandSurfaceJSON(&out); err != nil {
		t.Fatalf("writeCommandSurfaceJSON: %v", err)
	}
	var surface commandSurfaceSchema
	if err := json.Unmarshal(out.Bytes(), &surface); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if surface.SchemaVersion == "" {
		t.Fatal("schema_version is empty")
	}
	if !strings.HasPrefix(surface.SchemaVersion, "sdp-trace-command-surface-") {
		t.Fatalf("unexpected schema_version: %s", surface.SchemaVersion)
	}
}

func TestCommandSurfaceIncludesKnownProfiles(t *testing.T) {
	var out bytes.Buffer
	if err := writeCommandSurfaceJSON(&out); err != nil {
		t.Fatalf("writeCommandSurfaceJSON: %v", err)
	}
	var surface commandSurfaceSchema
	if err := json.Unmarshal(out.Bytes(), &surface); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	wantProfiles := map[string]bool{
		"adapter-capture":         true,
		"managed-harness":         true,
		"forensic-retention":      true,
		"ci-artifact-observation": true,
		"authority-envelope":      true,
	}
	got := map[string]bool{}
	for _, p := range surface.Profiles {
		got[p.ID] = true
	}
	for id := range wantProfiles {
		if !got[id] {
			t.Fatalf("missing profile %q", id)
		}
	}
}

func TestCommandSurfaceIncludesKnownWitnessKinds(t *testing.T) {
	var out bytes.Buffer
	if err := writeCommandSurfaceJSON(&out); err != nil {
		t.Fatalf("writeCommandSurfaceJSON: %v", err)
	}
	var surface commandSurfaceSchema
	if err := json.Unmarshal(out.Bytes(), &surface); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	want := map[string]bool{
		"github-actions": true,
		"gitlab-ci":    true,
		"buildkite":    true,
		"customer-pki": true,
	}
	got := map[string]bool{}
	for _, k := range surface.WitnessKinds {
		got[k] = true
	}
	for k := range want {
		if !got[k] {
			t.Fatalf("missing witness kind %q", k)
		}
	}
}

func TestCommandSurfaceIncludesKnownStates(t *testing.T) {
	var out bytes.Buffer
	if err := writeCommandSurfaceJSON(&out); err != nil {
		t.Fatalf("writeCommandSurfaceJSON: %v", err)
	}
	var surface commandSurfaceSchema
	if err := json.Unmarshal(out.Bytes(), &surface); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	want := map[string]bool{
		"observed":     true,
		"pass":         true,
		"fail":         true,
		"not_assessed": true,
		"cannot_verify": true,
	}
	got := map[string]bool{}
	for _, s := range surface.States {
		got[s.Name] = true
	}
	for k := range want {
		if !got[k] {
			t.Fatalf("missing state %q", k)
		}
	}
}

func TestCommandSurfaceUsageDrift(t *testing.T) {
	missing, stale, err := commandSurfaceUsageDrift()
	if err != nil {
		t.Fatalf("commandSurfaceUsageDrift: %v", err)
	}
	if len(missing) > 0 {
		t.Fatalf("registry usages missing from usageText: %s", strings.Join(missing, "; "))
	}
	if len(stale) > 0 {
		t.Fatalf("usageText usages stale vs registry: %s", strings.Join(stale, "; "))
	}
}

func TestRunCommandSurface(t *testing.T) {
	var out bytes.Buffer
	var errOut bytes.Buffer
	exit := run([]string{"command-surface"}, &out, &errOut)
	if exit != 0 {
		t.Fatalf("command-surface exit: %d err=%s", exit, errOut.String())
	}
	var surface commandSurfaceSchema
	if err := json.Unmarshal(out.Bytes(), &surface); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if surface.SchemaVersion == "" {
		t.Fatal("schema_version is empty")
	}
	if len(surface.Commands) == 0 {
		t.Fatal("no commands in surface")
	}
}

func TestRunCommandSurfaceRejectsExtraArgs(t *testing.T) {
	var out bytes.Buffer
	var errOut bytes.Buffer
	exit := run([]string{"command-surface", "unexpected"}, &out, &errOut)
	if exit != exitUsage {
		t.Fatalf("command-surface with extra args exit=%d, want exitUsage (%d); stderr=%q", exit, exitUsage, errOut.String())
	}
}

func TestCommandSurfaceDriftError(t *testing.T) {
	if err := commandSurfaceDriftError(nil, nil); err != nil {
		t.Fatalf("expected nil for empty drift, got %v", err)
	}
	if err := commandSurfaceDriftError([]string{"a"}, nil); err == nil {
		t.Fatal("expected error for missing drift")
	}
	if err := commandSurfaceDriftError(nil, []string{"b"}); err == nil {
		t.Fatal("expected error for stale drift")
	}
}

func jsonEqual(a, b commandSurfaceSchema) bool {
	// Simple structural equality via JSON round-trip for test determinism.
	ja, _ := json.Marshal(a)
	jb, _ := json.Marshal(b)
	return string(ja) == string(jb)
}
