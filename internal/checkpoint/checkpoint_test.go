package checkpoint

import (
	"context"
	"encoding/json"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/fall_out_bug/sdp-trace/internal/recorder"
	"github.com/fall_out_bug/sdp-trace/internal/trace"
)

func TestCreateAndVerifyLocalSignedCheckpoint(t *testing.T) {
	runDir := recordRun(t, "checkpoint-run", "task-1")
	key := GenerateKeyForTest(t, "local-dev")

	cp, err := Create(runDir, CreateOptions{
		CheckpointID: "checkpoint-001",
		Sequence:     0,
		SignerID:     "local-dev",
		Key:          key,
	})
	if err != nil {
		t.Fatalf("create checkpoint: %v", err)
	}

	result := Verify(runDir, cp, nil)
	if result.Result != StatePass {
		t.Fatalf("result = %s reasons=%v", result.Result, result.Reasons)
	}
	if result.TrustScope != TrustScopeLocalSigned {
		t.Fatalf("trust scope = %s", result.TrustScope)
	}
	if result.SignatureState != StatePass || result.ChainBindingState != StatePass || result.NonceBindingState != StatePass {
		t.Fatalf("unexpected states: %+v", result)
	}
	if result.SignerAuthorityState != StateNotAssessed {
		t.Fatalf("signer authority = %s", result.SignerAuthorityState)
	}
}

func TestVerifyFailsTamperedPayload(t *testing.T) {
	runDir := recordRun(t, "checkpoint-run", "task-1")
	key := GenerateKeyForTest(t, "local-dev")
	cp, err := Create(runDir, CreateOptions{CheckpointID: "checkpoint-001", SignerID: "local-dev", Key: key})
	if err != nil {
		t.Fatal(err)
	}
	cp.Payload.TaskHash = trace.SHA256Hex("different-task")

	result := Verify(runDir, cp, nil)
	if result.Result != StateFail {
		t.Fatalf("result = %s reasons=%v", result.Result, result.Reasons)
	}
	if result.PayloadDigestState != StateFail || result.SignatureState != StateFail {
		t.Fatalf("expected payload and signature failure, got %+v", result)
	}
}

func TestVerifyFailsWrongSchemaVersion(t *testing.T) {
	runDir := recordRun(t, "checkpoint-run", "task-1")
	key := GenerateKeyForTest(t, "local-dev")
	cp, err := Create(runDir, CreateOptions{CheckpointID: "checkpoint-001", SignerID: "local-dev", Key: key})
	if err != nil {
		t.Fatal(err)
	}
	cp.SchemaVersion = "old-checkpoint-v0"

	result := Verify(runDir, cp, nil)
	if result.Result != StateFail {
		t.Fatalf("result = %s reasons=%v", result.Result, result.Reasons)
	}
	if result.TrustScope != TrustScopeUntrustedShape {
		t.Fatalf("trust scope = %s", result.TrustScope)
	}
}

func TestVerifyFailsReplayAgainstDifferentRun(t *testing.T) {
	runDir := recordRun(t, "checkpoint-run", "task-1")
	otherRunDir := recordRun(t, "other-run", "task-2")
	key := GenerateKeyForTest(t, "local-dev")
	cp, err := Create(runDir, CreateOptions{CheckpointID: "checkpoint-001", SignerID: "local-dev", Key: key})
	if err != nil {
		t.Fatal(err)
	}

	result := Verify(otherRunDir, cp, nil)
	if result.Result != StateFail {
		t.Fatalf("result = %s reasons=%v", result.Result, result.Reasons)
	}
	if result.RunBindingState != StateFail || result.SourceBindingState != StateFail || result.NonceBindingState != StateFail {
		t.Fatalf("expected replay binding failures, got %+v", result)
	}
}

func TestVerifyCheckpointSetFailsSequenceGapsAndDuplicates(t *testing.T) {
	runDir := recordRun(t, "checkpoint-run", "task-1")
	key := GenerateKeyForTest(t, "local-dev")
	first, err := Create(runDir, CreateOptions{CheckpointID: "checkpoint-000", Sequence: 0, SignerID: "local-dev", Key: key})
	if err != nil {
		t.Fatal(err)
	}
	duplicate, err := Create(runDir, CreateOptions{CheckpointID: "checkpoint-000-copy", Sequence: 0, SignerID: "local-dev", Key: key})
	if err != nil {
		t.Fatal(err)
	}
	gapped, err := Create(runDir, CreateOptions{CheckpointID: "checkpoint-002", Sequence: 2, PreviousCheckpointDigest: first.PayloadDigest, SignerID: "local-dev", Key: key})
	if err != nil {
		t.Fatal(err)
	}

	if result := VerifySet(runDir, []SignedCheckpoint{first, duplicate}, nil); result.SequenceState != StateFail {
		t.Fatalf("duplicate sequence state = %s", result.SequenceState)
	}
	if result := VerifySet(runDir, []SignedCheckpoint{first, gapped}, nil); result.SequenceState != StateFail {
		t.Fatalf("gap sequence state = %s", result.SequenceState)
	}
	if result := VerifySet(runDir, []SignedCheckpoint{gapped, first}, nil); result.SequenceState != StateFail {
		t.Fatalf("descending sequence state = %s", result.SequenceState)
	}
}

func TestCreateRejectsSequenceWithoutPreviousDigest(t *testing.T) {
	runDir := recordRun(t, "checkpoint-run", "task-1")
	key := GenerateKeyForTest(t, "local-dev")

	_, err := Create(runDir, CreateOptions{
		CheckpointID: "checkpoint-001",
		Sequence:     1,
		SignerID:     "local-dev",
		Key:          key,
	})
	if err == nil {
		t.Fatalf("expected missing previous checkpoint digest error")
	}
}

func TestVerifyCannotVerifyWhenRunNonceMissing(t *testing.T) {
	runDir := recordRun(t, "checkpoint-run", "task-1")
	key := GenerateKeyForTest(t, "local-dev")
	cp, err := Create(runDir, CreateOptions{CheckpointID: "checkpoint-001", SignerID: "local-dev", Key: key})
	if err != nil {
		t.Fatal(err)
	}
	removeRunNonce(t, runDir)

	result := Verify(runDir, cp, nil)
	if result.Result != StateCannotVerify {
		t.Fatalf("result = %s reasons=%v", result.Result, result.Reasons)
	}
}

func TestVerifyCheckpointSetFailsForgedCheckpoint(t *testing.T) {
	runDir := recordRun(t, "checkpoint-run", "task-1")
	key := GenerateKeyForTest(t, "local-dev")
	first, err := Create(runDir, CreateOptions{CheckpointID: "checkpoint-000", Sequence: 0, SignerID: "local-dev", Key: key})
	if err != nil {
		t.Fatal(err)
	}
	forged := first
	forged.CheckpointID = "checkpoint-001"
	forged.Sequence = 1
	forged.Payload.PreviousCheckpointDigest = first.PayloadDigest
	forged.Payload.EventCount = 99

	result := VerifySet(runDir, []SignedCheckpoint{first, forged}, nil)
	if result.Result != StateFail {
		t.Fatalf("forged checkpoint set result = %s reasons=%v", result.Result, result.Reasons)
	}
	if result.SignatureState != StateFail || result.PayloadDigestState != StateFail {
		t.Fatalf("expected forged checkpoint signature/digest failure, got %+v", result)
	}
}

func TestPolicyMismatchFailsSignerAuthority(t *testing.T) {
	runDir := recordRun(t, "checkpoint-run", "task-1")
	key := GenerateKeyForTest(t, "local-dev")
	cp, err := Create(runDir, CreateOptions{CheckpointID: "checkpoint-001", SignerID: "local-dev", Key: key})
	if err != nil {
		t.Fatal(err)
	}
	policy := TrustedCheckpointPolicy{
		SchemaVersion: PolicySchemaVersion,
		PolicyID:      "test-policy",
		AllowedSigners: []TrustedSigner{
			{SignerID: "ci-signer", Authority: AuthorityCIIsolatedJob},
		},
	}

	result := Verify(runDir, cp, &policy)
	if result.Result != StateFail {
		t.Fatalf("result = %s reasons=%v", result.Result, result.Reasons)
	}
	if result.SignerAuthorityState != StateFail {
		t.Fatalf("signer authority = %s", result.SignerAuthorityState)
	}
}

func TestPolicyWithoutPublicKeyCannotVerifySignerAuthority(t *testing.T) {
	runDir := recordRun(t, "checkpoint-run", "task-1")
	key := GenerateKeyForTest(t, "local-dev")
	cp, err := Create(runDir, CreateOptions{CheckpointID: "checkpoint-001", SignerID: "local-dev", Key: key})
	if err != nil {
		t.Fatal(err)
	}
	policy := TrustedCheckpointPolicy{
		SchemaVersion: PolicySchemaVersion,
		PolicyID:      "test-policy",
		AllowedSigners: []TrustedSigner{
			{SignerID: "local-dev", Authority: AuthorityLocalDevelopment},
		},
	}

	result := Verify(runDir, cp, &policy)
	if result.Result != StateCannotVerify {
		t.Fatalf("result = %s reasons=%v", result.Result, result.Reasons)
	}
	if result.SignerAuthorityState != StateCannotVerify {
		t.Fatalf("signer authority = %s", result.SignerAuthorityState)
	}
}

func TestPolicyPublicKeyMismatchFailsSignerAuthority(t *testing.T) {
	runDir := recordRun(t, "checkpoint-run", "task-1")
	key := GenerateKeyForTest(t, "local-dev")
	cp, err := Create(runDir, CreateOptions{CheckpointID: "checkpoint-001", SignerID: "local-dev", Key: key})
	if err != nil {
		t.Fatal(err)
	}
	otherKey := GenerateKeyForTest(t, "local-dev")
	policy := TrustedCheckpointPolicy{
		SchemaVersion: PolicySchemaVersion,
		PolicyID:      "test-policy",
		AllowedSigners: []TrustedSigner{
			{SignerID: "local-dev", Authority: AuthorityLocalDevelopment, PublicKey: otherKey.PublicKey},
		},
	}

	result := Verify(runDir, cp, &policy)
	if result.SignerAuthorityState != StateFail {
		t.Fatalf("signer authority = %s reasons=%v", result.SignerAuthorityState, result.Reasons)
	}
}

func TestPolicyAuthorityMismatchFailsSignerAuthority(t *testing.T) {
	runDir := recordRun(t, "checkpoint-run", "task-1")
	key := GenerateKeyForTest(t, "local-dev")
	cp, err := Create(runDir, CreateOptions{CheckpointID: "checkpoint-001", SignerID: "local-dev", Key: key})
	if err != nil {
		t.Fatal(err)
	}
	policy := TrustedCheckpointPolicy{
		SchemaVersion: PolicySchemaVersion,
		PolicyID:      "test-policy",
		AllowedSigners: []TrustedSigner{
			{SignerID: "local-dev", Authority: AuthorityCIIsolatedJob, PublicKey: cp.Signature.PublicKey},
		},
	}

	result := Verify(runDir, cp, &policy)
	if result.SignerAuthorityState != StateFail {
		t.Fatalf("signer authority = %s reasons=%v", result.SignerAuthorityState, result.Reasons)
	}
}

func TestPolicyAllowsLocalSignerAsPass(t *testing.T) {
	runDir := recordRun(t, "checkpoint-run", "task-1")
	key := GenerateKeyForTest(t, "local-dev")
	cp, err := Create(runDir, CreateOptions{CheckpointID: "checkpoint-001", SignerID: "local-dev", Key: key})
	if err != nil {
		t.Fatal(err)
	}
	policy := TrustedCheckpointPolicy{
		SchemaVersion: PolicySchemaVersion,
		PolicyID:      "test-policy",
		AllowedSigners: []TrustedSigner{
			{
				SignerID:  "local-dev",
				Authority: AuthorityLocalDevelopment,
				PublicKey: cp.Signature.PublicKey,
			},
		},
	}

	result := Verify(runDir, cp, &policy)
	if result.SignerAuthorityState != StatePass {
		t.Fatalf("signer authority = %s", result.SignerAuthorityState)
	}
	if result.TrustScope != TrustScopeLocalSigned {
		t.Fatalf("trust scope = %s", result.TrustScope)
	}
}

func TestPolicyUnknownAuthorityCannotVerifySignerAuthority(t *testing.T) {
	runDir := recordRun(t, "checkpoint-run", "task-1")
	key := GenerateKeyForTest(t, "local-dev")
	cp, err := Create(runDir, CreateOptions{CheckpointID: "checkpoint-001", SignerID: "local-dev", Key: key})
	if err != nil {
		t.Fatal(err)
	}
	cp.Signer.Authority = "mystery-authority"
	policy := TrustedCheckpointPolicy{
		SchemaVersion: PolicySchemaVersion,
		PolicyID:      "test-policy",
		AllowedSigners: []TrustedSigner{
			{
				SignerID:  "local-dev",
				Authority: "mystery-authority",
				PublicKey: cp.Signature.PublicKey,
			},
		},
	}

	result := Verify(runDir, cp, &policy)
	if result.SignerAuthorityState != StateCannotVerify {
		t.Fatalf("signer authority = %s", result.SignerAuthorityState)
	}
	if result.Result != StateCannotVerify {
		t.Fatalf("result = %s reasons=%v", result.Result, result.Reasons)
	}
}

func TestCheckpointJSONRoundTripDoesNotContainSensitiveCommandOutput(t *testing.T) {
	runDir := recordRunWithOutput(t, "checkpoint-run", "task-1", "SECRET_TOKEN_SHOULD_NOT_APPEAR")
	key := GenerateKeyForTest(t, "local-dev")
	cp, err := Create(runDir, CreateOptions{CheckpointID: "checkpoint-001", SignerID: "local-dev", Key: key})
	if err != nil {
		t.Fatal(err)
	}
	raw, err := json.MarshalIndent(cp, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(raw), "SECRET_TOKEN_SHOULD_NOT_APPEAR") {
		t.Fatalf("checkpoint leaked command output: %s", string(raw))
	}
}

func recordRun(t *testing.T, wrapperName, task string) string {
	t.Helper()
	return recordRunWithOutput(t, wrapperName, task, "ok")
}

func recordRunWithOutput(t *testing.T, wrapperName, task, output string) string {
	t.Helper()
	sh := mustFindShell(t)
	runDir := filepath.Join(t.TempDir(), "run")
	result, err := recorder.Run(context.Background(), recorder.RecorderOptions{
		WrapperName:        wrapperName,
		Task:               task,
		OutputDir:          runDir,
		UseDefaultContract: true,
		Command:            []string{sh, "-c", "printf '%s\\n' \"$1\"", "sh", output},
	})
	if err != nil {
		t.Fatalf("record run: %v", err)
	}
	if result.ExitCode != 0 {
		t.Fatalf("exit code = %d", result.ExitCode)
	}
	return runDir
}

func mustFindShell(t *testing.T) string {
	t.Helper()
	sh, err := exec.LookPath("sh")
	if err != nil {
		t.Skip("sh unavailable")
	}
	return sh
}

func GenerateKeyForTest(t *testing.T, signerID string) KeyPair {
	t.Helper()
	key, err := GenerateKeyPair(signerID)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	return key
}

func removeRunNonce(t *testing.T, runDir string) {
	t.Helper()
	artifact, err := trace.OpenRunArtifact(runDir)
	if err != nil {
		t.Fatal(err)
	}
	prev := trace.NullEventHash
	for index, event := range artifact.Events {
		if index == 0 {
			delete(event.EventPayload, "run_nonce")
		}
		event.PrevEventHash = prev
		computed, err := event.WithComputedEventHash()
		if err != nil {
			t.Fatal(err)
		}
		artifact.Events[index] = computed
		prev = computed.EventHash
		if err := artifact.Layout.WriteEvent(computed); err != nil {
			t.Fatal(err)
		}
	}
	artifact.Manifest.EventChainHead = prev
	artifact.Manifest.FinalChainHead = prev
	if err := artifact.Layout.WriteRun(artifact.Manifest); err != nil {
		t.Fatal(err)
	}
}
