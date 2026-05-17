package witness

import (
	"bytes"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"math/big"
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

func TestBuildProfileRejectsUnsupportedKind(t *testing.T) {
	if _, err := BuildProfile("unsupported", t.TempDir(), "", ProfileOptions{}); err == nil || !strings.Contains(err.Error(), "unsupported witness kind") {
		t.Fatalf("BuildProfile unsupported error = %v", err)
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

func TestBuildCIEnvelopeProfileIncludesReportArtifacts(t *testing.T) {
	root := writeRunRoot(t)
	reportDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(reportDir, "report.txt"), []byte("ok"), 0o644); err != nil {
		t.Fatalf("write report artifact: %v", err)
	}

	artifacts, err := hashRunArtifacts(root)
	if err != nil {
		t.Fatalf("hash artifacts: %v", err)
	}
	envelope := EnvelopeInput{
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
	envelopePath := writeJSON(t, t.TempDir(), "gitlab-envelope.json", envelope)

	record, err := BuildCIEnvelopeProfile(KindGitLabCI, root, reportDir, envelopePath)
	if err != nil {
		t.Fatalf("BuildCIEnvelopeProfile: %v", err)
	}
	if record.Status != StatusPass || record.Reason != ReasonProfileVerified {
		t.Fatalf("record = %s/%s", record.Status, record.Reason)
	}
	if len(record.ReportArtifacts) == 0 {
		t.Fatalf("report artifacts were not hashed")
	}
}

func TestValidateCIEnvelopeStates(t *testing.T) {
	pass := ProfileStates{
		IdentityState:        statePass,
		SignerAuthorityState: statePass,
		FreshnessState:       statePass,
		ArtifactBindingState: statePass,
		SourceBindingState:   statePass,
		RunBindingState:      statePass,
		PolicyBindingState:   statePass,
		IndependenceState:    independenceCIJob,
	}
	tests := []struct {
		name   string
		mutate func(*ProfileStates)
		status string
		scope  string
		reason string
	}{
		{
			name:   "pass",
			mutate: func(*ProfileStates) {},
		},
		{
			name:   "missing identity",
			mutate: func(state *ProfileStates) { state.IdentityState = stateCannotVerify },
			status: StatusCannotVerify,
			scope:  stateCannotVerify,
			reason: ReasonMissingIdentity,
		},
		{
			name:   "missing signer",
			mutate: func(state *ProfileStates) { state.SignerAuthorityState = stateCannotVerify },
			status: StatusCannotVerify,
			scope:  stateCannotVerify,
			reason: ReasonMissingSigner,
		},
		{
			name:   "stale freshness",
			mutate: func(state *ProfileStates) { state.FreshnessState = stateFail },
			status: StatusFail,
			scope:  stateFail,
			reason: ReasonStaleFreshness,
		},
		{
			name:   "missing freshness",
			mutate: func(state *ProfileStates) { state.FreshnessState = stateCannotVerify },
			status: StatusCannotVerify,
			scope:  stateCannotVerify,
			reason: ReasonMissingFreshness,
		},
		{
			name:   "source mismatch",
			mutate: func(state *ProfileStates) { state.SourceBindingState = stateFail },
			status: StatusFail,
			scope:  stateFail,
			reason: ReasonSourceMismatch,
		},
		{
			name:   "run mismatch",
			mutate: func(state *ProfileStates) { state.RunBindingState = stateFail },
			status: StatusFail,
			scope:  stateFail,
			reason: ReasonRunMismatch,
		},
		{
			name:   "missing policy",
			mutate: func(state *ProfileStates) { state.PolicyBindingState = stateCannotVerify },
			status: StatusCannotVerify,
			scope:  stateCannotVerify,
			reason: ReasonPolicyMissing,
		},
		{
			name:   "missing artifact",
			mutate: func(state *ProfileStates) { state.ArtifactBindingState = stateCannotVerify },
			status: StatusCannotVerify,
			scope:  stateCannotVerify,
			reason: ReasonMissingArtifact,
		},
		{
			name:   "environment-only",
			mutate: func(state *ProfileStates) { state.IndependenceState = "ci_same_job" },
			status: StatusCannotVerify,
			scope:  stateCannotVerify,
			reason: ReasonEnvOnly,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			state := pass
			tt.mutate(&state)
			decision := validateCIEnvelopeStates(state)
			if decision.status != tt.status || decision.scope != tt.scope || decision.reason != tt.reason {
				t.Fatalf("decision = %+v", decision)
			}
		})
	}
}

func TestBuildkiteEnvelopeSignerAndRunBindingFailures(t *testing.T) {
	root := writeRunRoot(t)
	artifacts, err := hashRunArtifacts(root)
	if err != nil {
		t.Fatalf("hash artifacts: %v", err)
	}
	base := EnvelopeInput{
		ProfileID:    "buildkite-v1",
		ProviderKind: KindBuildkite,
		Source:       SourceIdentity{Repository: "org/repo", Ref: "refs/heads/main", CommitSHA: "abc123"},
		CI:           CIIdentity{Provider: KindBuildkite, RunID: "pipeline-42", Job: "verify"},
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
			name: "missing signer",
			mutate: func(input *EnvelopeInput) {
				input.ProfileStates.SignerAuthorityState = stateCannotVerify
			},
			status: StatusCannotVerify,
			reason: ReasonMissingSigner,
		},
		{
			name: "run mismatch",
			mutate: func(input *EnvelopeInput) {
				input.ProfileStates.RunBindingState = stateFail
			},
			status: StatusFail,
			reason: ReasonRunMismatch,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := base
			tt.mutate(&input)
			envelopePath := writeJSON(t, t.TempDir(), "buildkite-envelope.json", input)
			record, err := BuildCIEnvelopeProfile(KindBuildkite, root, "", envelopePath)
			if err != nil {
				t.Fatalf("BuildCIEnvelopeProfile: %v", err)
			}
			if record.Status != tt.status || record.Reason != tt.reason {
				t.Fatalf("record = %s reason=%s", record.Status, record.Reason)
			}
		})
	}
}

func TestValidateCustomerPKIAuthority(t *testing.T) {
	publicKey, _, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	policyDigest := "cccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc"
	basePolicy := CustomerPKIAuthorityPolicy{
		SchemaVersion:   "sdp-trace-customer-pki-authority/v1",
		ProfileID:       "customer-pki-v1",
		AllowedSignerID: "signer-1",
		PublicKeySHA256: digestBytes(publicKey),
		PolicyDigest:    policyDigest,
		KeyCustodyState: "hsm",
	}
	baseFreshness := CustomerPKIFreshnessEvidence{
		SignerID:     "signer-1",
		PolicyDigest: policyDigest,
	}

	tests := []struct {
		name       string
		mutate     func(*CustomerPKIAuthorityPolicy)
		status     string
		reason     string
		stateField string
		stateValue string
	}{
		{
			name: "supports pass",
		},
		{
			name:       "unsupported profile",
			mutate:     func(policy *CustomerPKIAuthorityPolicy) { policy.ProfileID = "customer-pki-v2" },
			status:     StatusFail,
			reason:     ReasonSignerMismatch,
			stateField: "signer",
			stateValue: stateFail,
		},
		{
			name:       "empty signer id",
			mutate:     func(policy *CustomerPKIAuthorityPolicy) { policy.AllowedSignerID = "" },
			status:     StatusFail,
			reason:     ReasonSignerMismatch,
			stateField: "signer",
			stateValue: stateFail,
		},
		{
			name:       "signer mismatch",
			mutate:     func(policy *CustomerPKIAuthorityPolicy) { policy.AllowedSignerID = "other" },
			status:     StatusFail,
			reason:     ReasonSignerMismatch,
			stateField: "signer",
			stateValue: stateFail,
		},
		{
			name:       "public key mismatch",
			mutate:     func(policy *CustomerPKIAuthorityPolicy) { policy.PublicKeySHA256 = strings.Repeat("a", 64) },
			status:     StatusFail,
			reason:     ReasonSignerMismatch,
			stateField: "signer",
			stateValue: stateFail,
		},
		{
			name:       "policy mismatch",
			mutate:     func(policy *CustomerPKIAuthorityPolicy) { policy.PolicyDigest = strings.Repeat("b", 64) },
			status:     StatusFail,
			reason:     ReasonPolicyMismatch,
			stateField: "policy",
			stateValue: stateFail,
		},
		{
			name: "revocation required but missing state",
			mutate: func(policy *CustomerPKIAuthorityPolicy) {
				policy.RevocationRequired = true
				policy.RevocationState = ""
			},
			status:     StatusNotAssessed,
			reason:     ReasonRevocationNA,
			stateField: "signer",
			stateValue: stateNotAssessed,
		},
		{
			name:       "revoked",
			mutate:     func(policy *CustomerPKIAuthorityPolicy) { policy.RevocationState = "revoked" },
			status:     StatusFail,
			reason:     ReasonCertRevoked,
			stateField: "signer",
			stateValue: stateFail,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			policy := basePolicy
			if tt.mutate != nil {
				tt.mutate(&policy)
			}
			record := baseRecord(KindCustomerPKI)
			states := customerPKIPassStates(policy)
			ok := validateCustomerPKIAuthority(&record, states, publicKey, policy, baseFreshness)
			if tt.status == "" {
				if !ok {
					t.Fatalf("expected authority validation pass")
				}
				if states.SignerAuthorityState != statePass {
					t.Fatalf("signer authority state = %s", states.SignerAuthorityState)
				}
				return
			}
			if ok {
				t.Fatalf("expected authority validation fail")
			}
			if record.Status != tt.status {
				t.Fatalf("status = %s", record.Status)
			}
			if record.Reason != tt.reason {
				t.Fatalf("reason = %s", record.Reason)
			}
			if tt.stateField == "signer" && states.SignerAuthorityState != tt.stateValue {
				t.Fatalf("signer authority state = %s", states.SignerAuthorityState)
			}
			if tt.stateField == "policy" && states.PolicyBindingState != tt.stateValue {
				t.Fatalf("policy binding state = %s", states.PolicyBindingState)
			}
		})
	}
}

func TestStrongDigestAcceptsSHA256OrStrongerHex(t *testing.T) {
	for _, digest := range []string{
		strings.Repeat("a", 64),
		strings.Repeat("b", 96),
		strings.Repeat("c", 128),
	} {
		if !strongDigest(digest) {
			t.Fatalf("digest should be accepted as sha256-or-stronger: len=%d", len(digest))
		}
	}
	if strongDigest(strings.Repeat("d", 63)) {
		t.Fatalf("short digest accepted")
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
	if record.ProfileStates.SourceBindingState != stateNotAssessed {
		t.Fatalf("source binding = %s", record.ProfileStates.SourceBindingState)
	}
}

func TestLoadCustomerPublicKeyFromPublicPEM(t *testing.T) {
	publicKey, _, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	dir := t.TempDir()
	publicKeyPath := writePublicKey(t, dir, publicKey)
	loaded, err := loadCustomerPublicKey(ProfileOptions{
		CustomerPKIPublicKey: publicKeyPath,
	})
	if err != nil {
		t.Fatalf("loadCustomerPublicKey: %v", err)
	}
	if !bytes.Equal(loaded, publicKey) {
		t.Fatalf("public key mismatch")
	}
}

func TestLoadCustomerPublicKeyFromX509Certificate(t *testing.T) {
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	certPath := writeX509Certificate(t, t.TempDir(), publicKey, privateKey)
	loaded, err := loadCustomerPublicKey(ProfileOptions{
		CustomerPKIPublicCert: certPath,
	})
	if err != nil {
		t.Fatalf("loadCustomerPublicKey: %v", err)
	}
	if !bytes.Equal(loaded, publicKey) {
		t.Fatalf("public key mismatch")
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
		key      ed25519.PublicKey
		reason   string
		status   string
	}{
		{
			name: "unsupported profile",
			policy: func() CustomerPKIAuthorityPolicy {
				copy := policy
				copy.ProfileID = "customer-pki-v2"
				return copy
			}(),
			evidence: evidence,
			key:      publicKey,
			status:   StatusFail,
			reason:   ReasonSignerMismatch,
		},
		{
			name: "empty signer",
			policy: func() CustomerPKIAuthorityPolicy {
				copy := policy
				copy.AllowedSignerID = ""
				return copy
			}(),
			evidence: evidence,
			key:      publicKey,
			status:   StatusFail,
			reason:   ReasonSignerMismatch,
		},
		{
			name: "signer mismatch",
			policy: func() CustomerPKIAuthorityPolicy {
				copy := policy
				copy.AllowedSignerID = "other-signer"
				return copy
			}(),
			evidence: evidence,
			key:      publicKey,
			status:   StatusFail,
			reason:   ReasonSignerMismatch,
		},
		{
			name: "public key mismatch",
			policy: func() CustomerPKIAuthorityPolicy {
				copy := policy
				copy.PublicKeySHA256 = strings.Repeat("c", 64)
				return copy
			}(),
			evidence: evidence,
			key:      publicKey,
			status:   StatusFail,
			reason:   ReasonSignerMismatch,
		},
		{
			name: "policy digest mismatch",
			policy: func() CustomerPKIAuthorityPolicy {
				copy := policy
				copy.PolicyDigest = strings.Repeat("c", 64)
				return copy
			}(),
			evidence: evidence,
			key:      publicKey,
			status:   StatusFail,
			reason:   ReasonPolicyMismatch,
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
			key:    publicKey,
			status: StatusFail,
			reason: ReasonStaleFreshness,
		},
		{
			name:   "run mismatch",
			policy: policy,
			evidence: func() CustomerPKIFreshnessEvidence {
				copy := evidence
				copy.RunID = "other-run"
				copy.Signature = base64.StdEncoding.EncodeToString(ed25519.Sign(privateKey, []byte(freshnessPayload(copy))))
				return copy
			}(),
			key:    publicKey,
			status: StatusFail,
			reason: ReasonRunMismatch,
		},
		{
			name:   "invalid signature",
			policy: policy,
			evidence: func() CustomerPKIFreshnessEvidence {
				copy := evidence
				copy.Signature = base64.StdEncoding.EncodeToString(ed25519.Sign(privateKey, []byte("tampered")))
				return copy
			}(),
			key:    publicKey,
			status: StatusFail,
			reason: ReasonSignerMismatch,
		},
		{
			name: "revocation not assessed",
			policy: func() CustomerPKIAuthorityPolicy {
				copy := policy
				copy.RevocationRequired = true
				return copy
			}(),
			evidence: evidence,
			key:      publicKey,
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
			key:      publicKey,
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
			key:    publicKey,
			status: StatusFail,
			reason: ReasonArtifactMismatch,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			policyPath := writeJSON(t, dir, tt.name+"-policy.json", tt.policy)
			freshnessPath := writeJSON(t, dir, tt.name+"-freshness.json", tt.evidence)
			keyPath := publicKeyPath
			if tt.key != nil {
				keyPath = writePublicKey(t, dir, tt.key)
			}
			record, err := BuildCustomerPKI(root, "", ProfileOptions{
				CustomerPKIAuthorityPolicy: policyPath,
				CustomerPKIPublicKey:       keyPath,
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
	sentinel := "person@example.com"
	record, raw := writeGitLabWitnessProfileWithActor(t, "unsafe-serialized-envelope.json", sentinel)
	if record.Status != StatusFail || record.Reason != ReasonUnsafeOutput {
		t.Fatalf("record = %s reason=%s", record.Status, record.Reason)
	}
	if strings.Contains(string(raw), sentinel) {
		t.Fatalf("unsafe serialized sentinel leaked into witness output")
	}
}

func TestWriteProfileRejectsJWTShapedEnvelopeInput(t *testing.T) {
	sentinel := "eyJhbGciOiJFZERTQSJ9.eyJzdWIiOiJibG9jazIyIn0.signaturesecret"
	record, raw := writeGitLabWitnessProfileWithActor(t, "jwt-envelope.json", sentinel)
	if record.Status != StatusCannotVerify || record.Reason != ReasonMalformedInput {
		t.Fatalf("record = %s reason=%s", record.Status, record.Reason)
	}
	if strings.Contains(string(raw), sentinel) {
		t.Fatalf("jwt-shaped sentinel leaked into witness output")
	}
}

func writeGitLabWitnessProfileWithActor(t *testing.T, envelopeName, actor string) (Record, []byte) {
	t.Helper()
	root := writeRunRoot(t)
	dir := t.TempDir()
	artifacts, err := hashRunArtifacts(root)
	if err != nil {
		t.Fatalf("hash artifacts: %v", err)
	}
	envelope := EnvelopeInput{
		ProfileID:    "gitlab-ci-v1",
		ProviderKind: KindGitLabCI,
		Source:       SourceIdentity{Repository: "org/repo", Ref: "refs/heads/main", CommitSHA: "abc123"},
		CI:           CIIdentity{Provider: KindGitLabCI, RunID: "pipeline-42", Job: "verify", Actor: actor},
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
	envelopePath := writeJSON(t, dir, envelopeName, envelope)
	outPath := filepath.Join(dir, "witness.json")
	record, err := WriteProfile(KindGitLabCI, outPath, root, "", ProfileOptions{EnvelopePath: envelopePath})
	if err != nil {
		t.Fatalf("WriteProfile: %v", err)
	}
	raw, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("read witness: %v", err)
	}
	return record, raw
}

func TestFinalizeRecordRejectsJWTShapedSerializedOutput(t *testing.T) {
	sentinel := "eyJhbGciOiJFZERTQSJ9.eyJzdWIiOiJibG9jazIyIn0.signaturesecret"
	record := baseRecord(KindGitLabCI)
	record.ProfileID = sentinel
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

func writeX509Certificate(t *testing.T, dir string, publicKey ed25519.PublicKey, privateKey ed25519.PrivateKey) string {
	t.Helper()
	now := time.Now().UTC()
	template := x509.Certificate{
		SerialNumber:          big.NewInt(1),
		NotBefore:             now.Add(-time.Minute),
		NotAfter:              now.Add(time.Hour),
		KeyUsage:              x509.KeyUsageDigitalSignature,
		BasicConstraintsValid: true,
		SignatureAlgorithm:    x509.PureEd25519,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageCodeSigning},
	}
	der, err := x509.CreateCertificate(rand.Reader, &template, &template, publicKey, privateKey)
	if err != nil {
		t.Fatalf("create certificate: %v", err)
	}
	path := filepath.Join(dir, "public-cert.pem")
	raw := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der})
	if err := os.WriteFile(path, raw, 0o644); err != nil {
		t.Fatalf("write certificate: %v", err)
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
