package witness

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestGitLabEnvelopePassesWithBoundArtifacts(t *testing.T) {
	root := writeRunRoot(t)
	artifacts, err := hashRunArtifacts(root)
	if err != nil {
		t.Fatalf("hash artifacts: %v", err)
	}
	envelope := EnvelopeInput{
		SchemaVersion:       "sdp-trace-witness-profile-result/v1",
		ProfileID:           "gitlab-ci-v1",
		ProfileVersion:      "1.0",
		ProviderKind:        KindGitLabCI,
		RequestedTrustScope: TrustScopeCIWitnessed,
		Source:              SourceIdentity{Repository: "org/repo", Ref: "refs/heads/main", CommitSHA: "abc123"},
		CI:                  CIIdentity{Provider: KindGitLabCI, RunID: "pipeline-42", Job: "verify"},
		RunArtifacts:        artifacts,
		ProfileStates: ProfileStates{
			IdentityState:        statePass,
			SignerAuthorityState: statePass,
			FreshnessState:       statePass,
			ArtifactBindingState: statePass,
			SourceBindingState:   statePass,
			RunBindingState:      statePass,
			PolicyBindingState:   statePass,
			IndependenceState:    independenceCIJob,
		},
	}
	envelopePath := writeJSON(t, t.TempDir(), "gitlab-envelope.json", envelope)

	record, err := BuildCIEnvelopeProfile(KindGitLabCI, root, "", envelopePath)
	if err != nil {
		t.Fatalf("BuildCIEnvelopeProfile: %v", err)
	}
	if record.Status != StatusPass || record.TrustScope != TrustScopeCIWitnessed {
		t.Fatalf("record status/scope = %s/%s reason=%s", record.Status, record.TrustScope, record.Reason)
	}
	if record.ProfileStates.RunBindingState != statePass {
		t.Fatalf("run binding state = %s", record.ProfileStates.RunBindingState)
	}
}

func TestBuildkiteEnvelopePassesWithBoundArtifacts(t *testing.T) {
	root := writeRunRoot(t)
	artifacts, err := hashRunArtifacts(root)
	if err != nil {
		t.Fatalf("hash artifacts: %v", err)
	}
	envelope := EnvelopeInput{
		ProfileID:           "buildkite-v1",
		ProfileVersion:      "1.0",
		ProviderKind:        KindBuildkite,
		RequestedTrustScope: TrustScopeCIWitnessed,
		Source:              SourceIdentity{Repository: "org/repo", Ref: "refs/heads/main", CommitSHA: "abc123"},
		CI:                  CIIdentity{Provider: KindBuildkite, RunID: "pipeline-42", Job: "verify"},
		RunArtifacts:        artifacts,
		ProfileStates: ProfileStates{
			IdentityState:        statePass,
			SignerAuthorityState: statePass,
			FreshnessState:       statePass,
			ArtifactBindingState: statePass,
			SourceBindingState:   statePass,
			RunBindingState:      statePass,
			PolicyBindingState:   statePass,
			IndependenceState:    independenceCIJob,
		},
	}
	envelopePath := writeJSON(t, t.TempDir(), "buildkite-envelope.json", envelope)
	record, err := BuildCIEnvelopeProfile(KindBuildkite, root, "", envelopePath)
	if err != nil {
		t.Fatalf("BuildCIEnvelopeProfile: %v", err)
	}
	if record.Status != StatusPass || record.TrustScope != TrustScopeCIWitnessed {
		t.Fatalf("record status/scope = %s/%s reason=%s", record.Status, record.TrustScope, record.Reason)
	}
}

func TestBuildkiteWithoutEnvelopeCannotUpgradeFromEnvironment(t *testing.T) {
	clearProviderEnv(t)
	root := writeRunRoot(t)
	record, err := BuildCIEnvelopeProfile(KindBuildkite, root, "", "")
	if err != nil {
		t.Fatalf("BuildCIEnvelopeProfile: %v", err)
	}
	if record.Status != StatusCannotVerify {
		t.Fatalf("status = %s", record.Status)
	}
	if record.Reason != ReasonMissingIdentity {
		t.Fatalf("reason = %s", record.Reason)
	}
	if record.TrustScope == TrustScopeCIWitnessed || record.TrustScope == TrustScopeExternal {
		t.Fatalf("environment-only input upgraded trust: %s", record.TrustScope)
	}
}

func TestBuildkiteAmbientEnvironmentWithoutEnvelopeCannotUpgrade(t *testing.T) {
	clearProviderEnv(t)
	t.Setenv("BUILDKITE_BUILD_ID", "bk-build-1")
	root := writeRunRoot(t)
	record, err := BuildCIEnvelopeProfile(KindBuildkite, root, "", "")
	if err != nil {
		t.Fatalf("BuildCIEnvelopeProfile: %v", err)
	}
	if record.Status != StatusCannotVerify || record.Reason != ReasonEnvOnly {
		t.Fatalf("record = %s reason=%s", record.Status, record.Reason)
	}
	if record.TrustScope == TrustScopeCIWitnessed || record.TrustScope == TrustScopeExternal {
		t.Fatalf("ambient environment upgraded trust: %s", record.TrustScope)
	}
}

func TestGitLabEnvelopeArtifactMismatchFails(t *testing.T) {
	root := writeRunRoot(t)
	envelope := EnvelopeInput{
		ProfileID:    "gitlab-ci-v1",
		ProviderKind: KindGitLabCI,
		Source:       SourceIdentity{Repository: "org/repo", Ref: "refs/heads/main", CommitSHA: "abc123"},
		CI:           CIIdentity{Provider: KindGitLabCI, RunID: "pipeline-42", Job: "verify"},
		RunArtifacts: []ArtifactDigest{{Path: "001-agent-session/run.json", SHA256: "0000000000000000000000000000000000000000000000000000000000000000"}},
		ProfileStates: ProfileStates{
			IdentityState:        statePass,
			SignerAuthorityState: statePass,
			FreshnessState:       statePass,
			ArtifactBindingState: statePass,
			SourceBindingState:   statePass,
			RunBindingState:      statePass,
			PolicyBindingState:   statePass,
			IndependenceState:    independenceCIJob,
		},
	}
	envelopePath := writeJSON(t, t.TempDir(), "gitlab-envelope.json", envelope)

	record, err := BuildCIEnvelopeProfile(KindGitLabCI, root, "", envelopePath)
	if err != nil {
		t.Fatalf("BuildCIEnvelopeProfile: %v", err)
	}
	if record.Status != StatusFail || record.Reason != ReasonArtifactMismatch {
		t.Fatalf("record = %s reason=%s", record.Status, record.Reason)
	}
}

func TestCIEnvelopeNonPassReasonCodes(t *testing.T) {
	root := writeRunRoot(t)
	artifacts, err := hashRunArtifacts(root)
	if err != nil {
		t.Fatalf("hash artifacts: %v", err)
	}
	base := EnvelopeInput{
		ProfileID:    "gitlab-ci-v1",
		ProviderKind: KindGitLabCI,
		Source:       SourceIdentity{Repository: "org/repo", Ref: "refs/heads/main", CommitSHA: "abc123"},
		CI:           CIIdentity{Provider: KindGitLabCI, RunID: "pipeline-42", Job: "verify"},
		RunArtifacts: artifacts,
		ProfileStates: ProfileStates{
			IdentityState:        statePass,
			SignerAuthorityState: statePass,
			FreshnessState:       statePass,
			ArtifactBindingState: statePass,
			SourceBindingState:   statePass,
			RunBindingState:      statePass,
			PolicyBindingState:   statePass,
			IndependenceState:    independenceCIJob,
		},
	}
	tests := []struct {
		name   string
		mutate func(*EnvelopeInput)
		status string
		reason string
	}{
		{
			name: "missing identity",
			mutate: func(input *EnvelopeInput) {
				input.ProfileStates.IdentityState = stateCannotVerify
			},
			status: StatusCannotVerify,
			reason: ReasonMissingIdentity,
		},
		{
			name: "stale freshness",
			mutate: func(input *EnvelopeInput) {
				input.ProfileStates.FreshnessState = stateFail
			},
			status: StatusFail,
			reason: ReasonStaleFreshness,
		},
		{
			name: "run mismatch",
			mutate: func(input *EnvelopeInput) {
				input.ProfileStates.RunBindingState = stateFail
			},
			status: StatusFail,
			reason: ReasonRunMismatch,
		},
		{
			name: "source mismatch",
			mutate: func(input *EnvelopeInput) {
				input.ProfileStates.SourceBindingState = stateFail
			},
			status: StatusFail,
			reason: ReasonSourceMismatch,
		},
		{
			name: "same job topology cap",
			mutate: func(input *EnvelopeInput) {
				input.ProfileStates.IndependenceState = "ci_same_job"
			},
			status: StatusCannotVerify,
			reason: ReasonEnvOnly,
		},
		{
			name: "unsupported version",
			mutate: func(input *EnvelopeInput) {
				input.ProfileID = "gitlab-ci-v2"
			},
			status: StatusCannotVerify,
			reason: ReasonUnsupported,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := base
			tt.mutate(&input)
			envelopePath := writeJSON(t, t.TempDir(), "gitlab-envelope.json", input)
			record, err := BuildCIEnvelopeProfile(KindGitLabCI, root, "", envelopePath)
			if err != nil {
				t.Fatalf("BuildCIEnvelopeProfile: %v", err)
			}
			if record.Status != tt.status || record.Reason != tt.reason {
				t.Fatalf("record = %s reason=%s", record.Status, record.Reason)
			}
		})
	}
}

func TestCustomerPKIPassesWithSignedFreshnessEvidence(t *testing.T) {
	root := writeRunRootWithID(t, "run-1")
	dir := t.TempDir()
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	publicKeyPath := writePublicKey(t, dir, publicKey)
	payloadDigest := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	policyDigest := "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
	policy := CustomerPKIAuthorityPolicy{
		SchemaVersion:   "sdp-trace-customer-pki-authority/v1",
		ProfileID:       "customer-pki-v1",
		AllowedSignerID: "signer-1",
		PublicKeySHA256: digestBytes(publicKey),
		PolicyDigest:    policyDigest,
		KeyCustodyState: "hsm",
	}
	policyPath := writeJSON(t, dir, "policy.json", policy)
	evidence := CustomerPKIFreshnessEvidence{
		SchemaVersion: "sdp-trace-customer-pki-freshness/v1",
		SignerID:      "signer-1",
		PayloadDigest: payloadDigest,
		RunID:         "run-1",
		PolicyDigest:  policyDigest,
		IssuedAt:      time.Now().UTC().Add(-time.Minute).Format(time.RFC3339),
		ValidUntil:    time.Now().UTC().Add(time.Hour).Format(time.RFC3339),
		Nonce:         "nonce-1",
	}
	evidence.Signature = base64.StdEncoding.EncodeToString(ed25519.Sign(privateKey, []byte(freshnessPayload(evidence))))
	freshnessPath := writeJSON(t, dir, "freshness.json", evidence)

	record, err := BuildCustomerPKI(root, "", ProfileOptions{
		CustomerPKIAuthorityPolicy: policyPath,
		CustomerPKIPublicKey:       publicKeyPath,
		CustomerPKIPayloadDigest:   payloadDigest,
		CustomerPKIFreshness:       freshnessPath,
	})
	if err != nil {
		t.Fatalf("BuildCustomerPKI: %v", err)
	}
	if record.Status != StatusPass || record.TrustScope != TrustScopeExternal {
		t.Fatalf("record status/scope = %s/%s reason=%s", record.Status, record.TrustScope, record.Reason)
	}
	if record.ProfileStates.KeyCustodyState != "hsm" {
		t.Fatalf("key custody = %s", record.ProfileStates.KeyCustodyState)
	}
}

func TestCustomerPKIRejectsPrivateKeyInput(t *testing.T) {
	root := writeRunRootWithID(t, "run-1")
	dir := t.TempDir()
	privateKeyPath := filepath.Join(dir, "private.pem")
	if err := os.WriteFile(privateKeyPath, []byte("-----BEGIN PRIVATE KEY-----\nsecret\n-----END PRIVATE KEY-----\n"), 0o644); err != nil {
		t.Fatalf("write private key: %v", err)
	}
	record, err := BuildCustomerPKI(root, "", ProfileOptions{
		CustomerPKIAuthorityPolicy: filepath.Join(dir, "policy.json"),
		CustomerPKIPublicKey:       privateKeyPath,
		CustomerPKIPayloadDigest:   "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		CustomerPKIFreshness:       filepath.Join(dir, "freshness.json"),
	})
	if err != nil {
		t.Fatalf("BuildCustomerPKI: %v", err)
	}
	if record.Status != StatusFail || record.Reason != ReasonPrivateKeyInput {
		t.Fatalf("record = %s reason=%s", record.Status, record.Reason)
	}
}

func TestCustomerPKINonPassReasonCodes(t *testing.T) {
	root := writeRunRootWithID(t, "run-1")
	dir := t.TempDir()
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	publicKeyPath := writePublicKey(t, dir, publicKey)
	payloadDigest := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	policyDigest := "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
	policy := CustomerPKIAuthorityPolicy{
		SchemaVersion:   "sdp-trace-customer-pki-authority/v1",
		ProfileID:       "customer-pki-v1",
		AllowedSignerID: "signer-1",
		PublicKeySHA256: digestBytes(publicKey),
		PolicyDigest:    policyDigest,
		KeyCustodyState: "hsm",
	}
	evidence := CustomerPKIFreshnessEvidence{
		SchemaVersion: "sdp-trace-customer-pki-freshness/v1",
		SignerID:      "signer-1",
		PayloadDigest: payloadDigest,
		RunID:         "run-1",
		PolicyDigest:  policyDigest,
		IssuedAt:      time.Now().UTC().Add(-time.Minute).Format(time.RFC3339),
		ValidUntil:    time.Now().UTC().Add(time.Hour).Format(time.RFC3339),
		Nonce:         "nonce-1",
	}
	evidence.Signature = base64.StdEncoding.EncodeToString(ed25519.Sign(privateKey, []byte(freshnessPayload(evidence))))

	tests := []struct {
		name     string
		policy   CustomerPKIAuthorityPolicy
		evidence CustomerPKIFreshnessEvidence
		reason   string
		status   string
	}{
		{
			name: "signer mismatch",
			policy: func() CustomerPKIAuthorityPolicy {
				copy := policy
				copy.AllowedSignerID = "other-signer"
				return copy
			}(),
			evidence: evidence,
			status:   StatusFail,
			reason:   ReasonSignerMismatch,
		},
		{
			name:   "expired freshness",
			policy: policy,
			evidence: func() CustomerPKIFreshnessEvidence {
				copy := evidence
				copy.ValidUntil = time.Now().UTC().Add(-time.Hour).Format(time.RFC3339)
				copy.Signature = base64.StdEncoding.EncodeToString(ed25519.Sign(privateKey, []byte(freshnessPayload(copy))))
				return copy
			}(),
			status: StatusFail,
			reason: ReasonStaleFreshness,
		},
		{
			name: "revocation not assessed",
			policy: func() CustomerPKIAuthorityPolicy {
				copy := policy
				copy.RevocationRequired = true
				return copy
			}(),
			evidence: evidence,
			status:   StatusNotAssessed,
			reason:   ReasonRevocationNA,
		},
		{
			name: "revoked certificate",
			policy: func() CustomerPKIAuthorityPolicy {
				copy := policy
				copy.RevocationState = "revoked"
				return copy
			}(),
			evidence: evidence,
			status:   StatusFail,
			reason:   ReasonCertRevoked,
		},
		{
			name:   "weak digest",
			policy: policy,
			evidence: func() CustomerPKIFreshnessEvidence {
				copy := evidence
				copy.PayloadDigest = "weak"
				copy.Signature = base64.StdEncoding.EncodeToString(ed25519.Sign(privateKey, []byte(freshnessPayload(copy))))
				return copy
			}(),
			status: StatusFail,
			reason: ReasonArtifactMismatch,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			policyPath := writeJSON(t, dir, tt.name+"-policy.json", tt.policy)
			freshnessPath := writeJSON(t, dir, tt.name+"-freshness.json", tt.evidence)
			record, err := BuildCustomerPKI(root, "", ProfileOptions{
				CustomerPKIAuthorityPolicy: policyPath,
				CustomerPKIPublicKey:       publicKeyPath,
				CustomerPKIPayloadDigest:   payloadDigest,
				CustomerPKIFreshness:       freshnessPath,
			})
			if err != nil {
				t.Fatalf("BuildCustomerPKI: %v", err)
			}
			if record.Status != tt.status || record.Reason != tt.reason {
				t.Fatalf("record = %s reason=%s", record.Status, record.Reason)
			}
		})
	}
}

func TestUnsafeEnvelopeContentDoesNotLeakSecret(t *testing.T) {
	root := writeRunRoot(t)
	dir := t.TempDir()
	secret := "token_secret_block22_should_not_leak"
	envelopePath := filepath.Join(dir, "unsafe-envelope.json")
	if err := os.WriteFile(envelopePath, []byte(`{"schema_version":"sdp-trace-witness-profile-result/v1","raw_job_log":"`+secret+`"}`), 0o644); err != nil {
		t.Fatalf("write unsafe envelope: %v", err)
	}
	record, err := BuildCIEnvelopeProfile(KindGitLabCI, root, "", envelopePath)
	if err != nil {
		t.Fatalf("BuildCIEnvelopeProfile: %v", err)
	}
	raw, err := json.Marshal(record)
	if err != nil {
		t.Fatalf("marshal record: %v", err)
	}
	if strings.Contains(string(raw), secret) {
		t.Fatalf("unsafe sentinel leaked into witness output")
	}
	if record.Status != StatusCannotVerify || record.Reason != ReasonMalformedInput {
		t.Fatalf("record = %s reason=%s", record.Status, record.Reason)
	}
}

func TestSafeJSONRejectsSymlinkInput(t *testing.T) {
	dir := t.TempDir()
	targetPath := writeJSON(t, dir, "target.json", map[string]string{"ok": "true"})
	linkPath := filepath.Join(dir, "link.json")
	if err := os.Symlink(targetPath, linkPath); err != nil {
		t.Skipf("symlink unavailable: %v", err)
	}
	var decoded map[string]string
	if err := readSafeJSON(linkPath, &decoded); err == nil {
		t.Fatalf("expected symlink input to be rejected")
	}
}

func TestWriteProfileDetectsUnsafeSerializedOutput(t *testing.T) {
	root := writeRunRoot(t)
	dir := t.TempDir()
	artifacts, err := hashRunArtifacts(root)
	if err != nil {
		t.Fatalf("hash artifacts: %v", err)
	}
	sentinel := "person@example.com"
	envelope := EnvelopeInput{
		ProfileID:    "gitlab-ci-v1",
		ProviderKind: KindGitLabCI,
		Source:       SourceIdentity{Repository: "org/repo", Ref: "refs/heads/main", CommitSHA: "abc123"},
		CI:           CIIdentity{Provider: KindGitLabCI, RunID: "pipeline-42", Job: "verify", Actor: sentinel},
		RunArtifacts: artifacts,
		ProfileStates: ProfileStates{
			IdentityState:        statePass,
			SignerAuthorityState: statePass,
			FreshnessState:       statePass,
			ArtifactBindingState: statePass,
			SourceBindingState:   statePass,
			RunBindingState:      statePass,
			PolicyBindingState:   statePass,
			IndependenceState:    independenceCIJob,
		},
	}
	envelopePath := writeJSON(t, dir, "unsafe-serialized-envelope.json", envelope)
	outPath := filepath.Join(dir, "witness.json")
	record, err := WriteProfile(KindGitLabCI, outPath, root, "", ProfileOptions{EnvelopePath: envelopePath})
	if err != nil {
		t.Fatalf("WriteProfile: %v", err)
	}
	if record.Status != StatusFail || record.Reason != ReasonUnsafeOutput {
		t.Fatalf("record = %s reason=%s", record.Status, record.Reason)
	}
	raw, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("read witness: %v", err)
	}
	if strings.Contains(string(raw), sentinel) {
		t.Fatalf("unsafe serialized sentinel leaked into witness output")
	}
}

func TestWriteProfileRejectsJWTShapedEnvelopeInput(t *testing.T) {
	root := writeRunRoot(t)
	dir := t.TempDir()
	artifacts, err := hashRunArtifacts(root)
	if err != nil {
		t.Fatalf("hash artifacts: %v", err)
	}
	sentinel := "eyJhbGciOiJFZERTQSJ9.eyJzdWIiOiJibG9jazIyIn0.signaturesecret"
	envelope := EnvelopeInput{
		ProfileID:    "gitlab-ci-v1",
		ProviderKind: KindGitLabCI,
		Source:       SourceIdentity{Repository: "org/repo", Ref: "refs/heads/main", CommitSHA: "abc123"},
		CI:           CIIdentity{Provider: KindGitLabCI, RunID: "pipeline-42", Job: "verify", Actor: sentinel},
		RunArtifacts: artifacts,
		ProfileStates: ProfileStates{
			IdentityState:        statePass,
			SignerAuthorityState: statePass,
			FreshnessState:       statePass,
			ArtifactBindingState: statePass,
			SourceBindingState:   statePass,
			RunBindingState:      statePass,
			PolicyBindingState:   statePass,
			IndependenceState:    independenceCIJob,
		},
	}
	envelopePath := writeJSON(t, dir, "jwt-envelope.json", envelope)
	outPath := filepath.Join(dir, "witness.json")
	record, err := WriteProfile(KindGitLabCI, outPath, root, "", ProfileOptions{EnvelopePath: envelopePath})
	if err != nil {
		t.Fatalf("WriteProfile: %v", err)
	}
	if record.Status != StatusCannotVerify || record.Reason != ReasonMalformedInput {
		t.Fatalf("record = %s reason=%s", record.Status, record.Reason)
	}
	raw, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("read witness: %v", err)
	}
	if strings.Contains(string(raw), sentinel) {
		t.Fatalf("jwt-shaped sentinel leaked into witness output")
	}
}

func TestFinalizeRecordRejectsJWTShapedSerializedOutput(t *testing.T) {
	sentinel := "eyJhbGciOiJFZERTQSJ9.eyJzdWIiOiJibG9jazIyIn0.signaturesecret"
	record := baseRecord(KindGitLabCI)
	record.CI = CIIdentity{Provider: KindGitLabCI, RunID: "pipeline-42", Actor: sentinel}

	safe := finalizeRecordForWrite(record)

	if safe.Status != StatusFail || safe.Reason != ReasonUnsafeOutput {
		t.Fatalf("record = %s reason=%s", safe.Status, safe.Reason)
	}
	raw, err := json.Marshal(safe)
	if err != nil {
		t.Fatalf("marshal safe record: %v", err)
	}
	if strings.Contains(string(raw), sentinel) {
		t.Fatalf("jwt-shaped sentinel leaked into finalized witness output")
	}
}

func writeRunRoot(t *testing.T) string {
	return writeRunRootWithID(t, "pipeline-42")
}

func writeRunRootWithID(t *testing.T, runID string) string {
	t.Helper()
	root := t.TempDir()
	runDir := filepath.Join(root, "001-agent-session")
	if err := os.MkdirAll(runDir, 0o755); err != nil {
		t.Fatalf("mkdir run: %v", err)
	}
	if err := os.WriteFile(filepath.Join(runDir, "run.json"), []byte(`{"run_id":"`+runID+`"}`), 0o644); err != nil {
		t.Fatalf("write run: %v", err)
	}
	return root
}

func writeJSON(t *testing.T, dir, name string, value any) string {
	t.Helper()
	raw, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		t.Fatalf("marshal %s: %v", name, err)
	}
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, append(raw, '\n'), 0o644); err != nil {
		t.Fatalf("write %s: %v", name, err)
	}
	return path
}

func writePublicKey(t *testing.T, dir string, publicKey ed25519.PublicKey) string {
	t.Helper()
	der, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		t.Fatalf("marshal public key: %v", err)
	}
	path := filepath.Join(dir, "public.pem")
	raw := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: der})
	if err := os.WriteFile(path, raw, 0o644); err != nil {
		t.Fatalf("write public key: %v", err)
	}
	return path
}

func clearProviderEnv(t *testing.T) {
	t.Helper()
	for _, key := range []string{
		"BUILDKITE",
		"BUILDKITE_BUILD_ID",
		"BUILDKITE_JOB_ID",
		"BUILDKITE_COMMIT",
		"GITLAB_CI",
		"CI_PIPELINE_ID",
		"CI_JOB_ID",
		"CI_COMMIT_SHA",
	} {
		t.Setenv(key, "")
	}
}
