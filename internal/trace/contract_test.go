package trace

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestLoadContractDefaultsAndErrors(t *testing.T) {
	if contract, err := LoadContract(""); err != nil || contract.ContractID != DefaultContract.ContractID {
		t.Fatalf("empty LoadContract() = %+v err=%v", contract, err)
	}

	dir := t.TempDir()
	path := filepath.Join(dir, "custom-contract.json")
	if err := os.WriteFile(path, []byte(`{"contract_id":"custom","required_events":["run_started"]}`), 0o644); err != nil {
		t.Fatalf("write contract: %v", err)
	}
	contract, err := LoadContract(path)
	if err != nil {
		t.Fatalf("LoadContract() error = %v", err)
	}
	if contract.ContractID != "custom" || contract.Version != SchemaVersion || len(contract.RequiredEvents) != 1 {
		t.Fatalf("contract defaults = %+v", contract)
	}

	if err := os.WriteFile(path, []byte(`{"contract_id":`), 0o644); err != nil {
		t.Fatalf("write malformed contract: %v", err)
	}
	if _, err := LoadContract(path); err == nil {
		t.Fatalf("expected malformed contract error")
	}
}

func TestContractWithDefaultsUsesFileNameAndDefaultEvents(t *testing.T) {
	contract := (Contract{}).withDefaults(filepath.Join("contracts", "baseline.json"))
	if contract.ContractID != "baseline.json" || contract.Version != SchemaVersion {
		t.Fatalf("identity defaults = %+v", contract)
	}
	if len(contract.RequiredEvents) != len(DefaultContract.RequiredEvents) {
		t.Fatalf("required events = %+v", contract.RequiredEvents)
	}
	contract.RequiredEvents[0] = "mutated"
	if DefaultContract.RequiredEvents[0] == "mutated" {
		t.Fatalf("defaults were aliased")
	}
}

func TestLocalSourceSnapshotReportsGitCleanDirtyAndFallback(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available")
	}
	repo := t.TempDir()
	runTraceGit(t, repo, "init")
	runTraceGit(t, repo, "config", "user.email", "trace@example.com")
	runTraceGit(t, repo, "config", "user.name", "Trace Test")
	if err := os.WriteFile(filepath.Join(repo, "file.txt"), []byte("one\n"), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}
	runTraceGit(t, repo, "add", "file.txt")
	runTraceGit(t, repo, "commit", "-m", "initial")

	digest, state := LocalSourceSnapshot(repo)
	if digest == "" || state != "git_tree_clean" {
		t.Fatalf("clean snapshot digest=%q state=%q", digest, state)
	}
	if err := os.WriteFile(filepath.Join(repo, "file.txt"), []byte("two\n"), 0o644); err != nil {
		t.Fatalf("dirty file: %v", err)
	}
	_, state = LocalSourceSnapshot(repo)
	if state != "git_tree_dirty" {
		t.Fatalf("dirty state = %q", state)
	}

	digest, state = LocalSourceSnapshot(filepath.Join(t.TempDir(), "not-a-repo"))
	if digest == "" || state != "not_assessed" {
		t.Fatalf("fallback snapshot digest=%q state=%q", digest, state)
	}
}

func TestMissingEvidenceAndPathHelpers(t *testing.T) {
	table := GenerateMissingEvidenceTable(Contract{ContractID: "contract", RequiredEvents: []string{"a", "b"}}, map[string]bool{"a": true})
	if table.ContractID != "contract" || len(table.Rows) != 1 || table.Rows[0].ExpectedEvent != "b" {
		t.Fatalf("missing evidence table = %+v", table)
	}
	base := t.TempDir()
	if got := ResolveContractPath(base, "contract.json"); got != filepath.Join(base, "contract.json") {
		t.Fatalf("relative contract path = %q", got)
	}
	if got := ResolveContractPath(base, ""); got != "" {
		t.Fatalf("empty contract path = %q", got)
	}
	abs := filepath.Join(base, "abs.json")
	if got := ResolveContractPath("/other", abs); got != abs {
		t.Fatalf("absolute contract path = %q", got)
	}
	if got := ResolveContractPath("", "contract.json"); got != "contract.json" {
		t.Fatalf("relative without base = %q", got)
	}
}

func TestJSONAndEventHashHelpers(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "payload.json")
	if err := os.WriteFile(path, []byte(`{"ok":true}`), 0o644); err != nil {
		t.Fatalf("write json: %v", err)
	}
	var decoded map[string]bool
	if err := ReadJSON(context.Background(), path, &decoded); err != nil {
		t.Fatalf("ReadJSON() error = %v", err)
	}
	if !decoded["ok"] {
		t.Fatalf("decoded = %+v", decoded)
	}
	if seq, err := EventSeqFromFilename("000123-run_started.json"); err != nil || seq != 123 {
		t.Fatalf("EventSeqFromFilename() = %d err=%v", seq, err)
	}
	if _, err := EventSeqFromFilename("run_started.json"); err == nil {
		t.Fatalf("expected invalid event filename error")
	}
	event := mustNewTestEvent(t, 0, EventRunStarted, NullEventHash)
	if hash, err := EventHash(event); err != nil || hash == "" {
		t.Fatalf("EventHash() = %q err=%v", hash, err)
	}
}

func TestEventPayloadDigestBranches(t *testing.T) {
	event := mustNewTestEvent(t, 0, EventRunStarted, NullEventHash)
	event.PayloadDigest = ""
	if err := event.VerifyPayloadDigest(); err != nil {
		t.Fatalf("empty digest should be skipped: %v", err)
	}

	event = mustNewTestEvent(t, 0, EventRunStarted, NullEventHash)
	event.EventPayload = nil
	event.Payload = []byte(`{"state":"run_started"}`)
	digest, err := CanonicalEventPayloadDigest([]byte(`{"state":"run_started"}`))
	if err != nil {
		t.Fatalf("payload digest: %v", err)
	}
	event.PayloadDigest = digest
	if err := event.VerifyPayloadDigest(); err != nil {
		t.Fatalf("decoded payload digest should verify: %v", err)
	}
	event.Payload = []byte(`{`)
	if err := event.VerifyPayloadDigest(); err == nil {
		t.Fatalf("expected invalid encoded payload error")
	}

	event = mustNewTestEvent(t, 0, EventRunStarted, NullEventHash)
	event.PayloadDigest = "sha256:wrong"
	if err := event.VerifyPayloadDigest(); err == nil {
		t.Fatalf("expected payload digest mismatch")
	}
}

func runTraceGit(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v failed: %v\n%s", args, err, out)
	}
}
