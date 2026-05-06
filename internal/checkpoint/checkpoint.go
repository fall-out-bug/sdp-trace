package checkpoint

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"github.com/fall_out_bug/sdp-trace/internal/trace"
)

const (
	CheckpointSchemaVersion     = "block15-signed-checkpoint-v1"
	VerificationSchemaVersion   = "block15-checkpoint-verification-v1"
	PolicySchemaVersion         = "block15-trusted-checkpoint-policy-v1"
	ProfileEd25519Detached      = "sdp-trace-checkpoint/ed25519-detached-v1"
	HashAlgorithmSHA256         = "sha256"
	SignatureAlgorithmEd25519   = "ed25519"
	AuthorityLocalDevelopment   = "local_development_key"
	AuthorityCIIsolatedJob      = "ci_isolated_job"
	AuthorityExternalWitness    = "external_witness_service"
	KeyIsolationNotAssessed     = "not_assessed"
	StatePass                   = "pass"
	StateFail                   = "fail"
	StateCannotVerify           = "cannot_verify"
	StateNotAssessed            = "not_assessed"
	StateNotIntegrated          = "not_integrated"
	TrustScopeLocalSigned       = "local_signed"
	TrustScopeCISigned          = "ci_signed"
	TrustScopeExternalWitnessed = "external_witnessed"
	TrustScopeUntrustedShape    = "untrusted_shape_only"
)

type KeyPair struct {
	Algorithm  string `json:"algorithm"`
	SignerID   string `json:"signer_id"`
	PublicKey  string `json:"public_key"`
	PrivateKey string `json:"private_key"`
}

type CreateOptions struct {
	CheckpointID             string
	Sequence                 int
	PreviousCheckpointDigest string
	SignerID                 string
	Key                      KeyPair
}

type SignedCheckpoint struct {
	SchemaVersion string                 `json:"schema_version"`
	CheckpointID  string                 `json:"checkpoint_id"`
	RunID         string                 `json:"run_id"`
	Sequence      int                    `json:"sequence"`
	Profile       string                 `json:"profile"`
	Canonical     trace.Canonicalization `json:"canonicalization"`
	HashAlgorithm string                 `json:"hash_algorithm"`
	Payload       Payload                `json:"payload"`
	PayloadDigest string                 `json:"payload_digest"`
	Signature     Signature              `json:"signature"`
	Signer        SignerIdentity         `json:"signer_identity"`
}

type Payload struct {
	RunID                    string        `json:"run_id"`
	RunNonce                 string        `json:"run_nonce"`
	EventChainHead           string        `json:"event_chain_head"`
	EventCount               int           `json:"event_count"`
	SourceSnapshotDigest     string        `json:"source_snapshot_digest"`
	SourceSnapshotState      string        `json:"source_snapshot_state"`
	TaskHash                 string        `json:"task_hash"`
	ContractDigest           string        `json:"contract_digest"`
	PreviousCheckpointDigest string        `json:"previous_checkpoint_digest"`
	ReplayContext            ReplayContext `json:"replay_context"`
}

type ReplayContext struct {
	Repository string `json:"repository"`
	Ref        string `json:"ref"`
	CommitSHA  string `json:"commit_sha"`
}

type Signature struct {
	Algorithm string `json:"algorithm"`
	Signature string `json:"signature"`
	PublicKey string `json:"public_key"`
}

type SignerIdentity struct {
	SignerID     string `json:"signer_id"`
	Authority    string `json:"authority"`
	KeyIsolation string `json:"key_isolation"`
}

type TrustedCheckpointPolicy struct {
	SchemaVersion  string          `json:"schema_version"`
	PolicyID       string          `json:"policy_id"`
	AllowedSigners []TrustedSigner `json:"allowed_signers"`
}

type TrustedSigner struct {
	SignerID  string `json:"signer_id"`
	Authority string `json:"authority"`
	PublicKey string `json:"public_key,omitempty"`
}

type VerificationResult struct {
	SchemaVersion        string   `json:"schema_version"`
	CheckpointID         string   `json:"checkpoint_id"`
	RunID                string   `json:"run_id"`
	Result               string   `json:"result"`
	TrustScope           string   `json:"trust_scope"`
	PayloadDigestState   string   `json:"payload_digest_state"`
	SignatureState       string   `json:"signature_state"`
	RunBindingState      string   `json:"run_binding_state"`
	ChainBindingState    string   `json:"chain_binding_state"`
	SourceBindingState   string   `json:"source_binding_state"`
	NonceBindingState    string   `json:"nonce_binding_state"`
	SequenceState        string   `json:"sequence_state"`
	SignerAuthorityState string   `json:"signer_authority_state"`
	ReplayFreshnessState string   `json:"replay_freshness_state"`
	Reasons              []string `json:"reasons"`
}

func GenerateKeyPair(signerID string) (KeyPair, error) {
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return KeyPair{}, err
	}
	return KeyPair{
		Algorithm:  SignatureAlgorithmEd25519,
		SignerID:   signerID,
		PublicKey:  base64.StdEncoding.EncodeToString(publicKey),
		PrivateKey: base64.StdEncoding.EncodeToString(privateKey),
	}, nil
}

func Create(runDir string, options CreateOptions) (SignedCheckpoint, error) {
	if strings.TrimSpace(options.CheckpointID) == "" {
		return SignedCheckpoint{}, errors.New("checkpoint_id is required")
	}
	if strings.TrimSpace(options.SignerID) == "" {
		return SignedCheckpoint{}, errors.New("signer_id is required")
	}
	if err := validateSequenceLink(options.Sequence, options.PreviousCheckpointDigest); err != nil {
		return SignedCheckpoint{}, err
	}
	privateKey, err := decodePrivateKey(options.Key)
	if err != nil {
		return SignedCheckpoint{}, err
	}
	if err := validateKeyPair(options.Key, privateKey); err != nil {
		return SignedCheckpoint{}, err
	}
	payload, err := BuildPayload(runDir, options.PreviousCheckpointDigest)
	if err != nil {
		return SignedCheckpoint{}, err
	}
	canonical, err := trace.CanonicalJSON(payload)
	if err != nil {
		return SignedCheckpoint{}, err
	}
	signature := ed25519.Sign(privateKey, canonical)
	publicKey := options.Key.PublicKey
	if publicKey == "" {
		publicKey = base64.StdEncoding.EncodeToString(privateKey.Public().(ed25519.PublicKey))
	}
	return SignedCheckpoint{
		SchemaVersion: CheckpointSchemaVersion,
		CheckpointID:  options.CheckpointID,
		RunID:         payload.RunID,
		Sequence:      options.Sequence,
		Profile:       ProfileEd25519Detached,
		Canonical: trace.Canonicalization{
			Algorithm: trace.CanonicalSchemaAlgo,
			Version:   trace.CanonicalAlgoVersion,
		},
		HashAlgorithm: HashAlgorithmSHA256,
		Payload:       payload,
		PayloadDigest: trace.SHA256Hex(string(canonical)),
		Signature: Signature{
			Algorithm: SignatureAlgorithmEd25519,
			Signature: base64.StdEncoding.EncodeToString(signature),
			PublicKey: publicKey,
		},
		Signer: SignerIdentity{
			SignerID:     options.SignerID,
			Authority:    AuthorityLocalDevelopment,
			KeyIsolation: KeyIsolationNotAssessed,
		},
	}, nil
}

func BuildPayload(runDir, previousCheckpointDigest string) (Payload, error) {
	artifact, err := trace.OpenRunArtifact(runDir)
	if err != nil {
		return Payload{}, err
	}
	if err := artifact.Manifest.Validate(); err != nil {
		return Payload{}, err
	}
	if artifact.Manifest.EventCount != len(artifact.Events) {
		return Payload{}, fmt.Errorf("event_count mismatch: run.json=%d files=%d", artifact.Manifest.EventCount, len(artifact.Events))
	}
	if err := trace.ValidateEventChain(artifact.Events); err != nil {
		return Payload{}, err
	}
	chainHead := artifact.Manifest.EventChainHead
	if chainHead == "" {
		chainHead = artifact.Manifest.FinalChainHead
	}
	nonce := runNonce(artifact.Events)
	if nonce == "" {
		return Payload{}, errors.New("run nonce missing from recorder_attached event")
	}
	return Payload{
		RunID:                    artifact.Manifest.RunID,
		RunNonce:                 nonce,
		EventChainHead:           chainHead,
		EventCount:               artifact.Manifest.EventCount,
		SourceSnapshotDigest:     artifact.Manifest.SourceSnapshot,
		SourceSnapshotState:      artifact.Manifest.SourceState,
		TaskHash:                 trace.SHA256Hex(artifact.Manifest.Task),
		ContractDigest:           artifact.Manifest.ContractDigest,
		PreviousCheckpointDigest: previousCheckpointDigest,
		ReplayContext: ReplayContext{
			Repository: "not_assessed",
			Ref:        "not_assessed",
			CommitSHA:  "not_assessed",
		},
	}, nil
}

func Verify(runDir string, checkpoint SignedCheckpoint, policy *TrustedCheckpointPolicy) VerificationResult {
	result := baseResult(checkpoint)
	if err := validateEnvelope(checkpoint); err != nil {
		result.Result = StateFail
		result.TrustScope = TrustScopeUntrustedShape
		result.Reasons = append(result.Reasons, err.Error())
		return result
	}
	if err := validateSequenceLink(checkpoint.Sequence, checkpoint.Payload.PreviousCheckpointDigest); err != nil {
		result.Result = StateFail
		result.TrustScope = TrustScopeUntrustedShape
		result.SequenceState = StateFail
		result.Reasons = append(result.Reasons, err.Error())
		return result
	}
	expected, err := BuildPayload(runDir, checkpoint.Payload.PreviousCheckpointDigest)
	if err != nil {
		result.Result = StateCannotVerify
		result.Reasons = append(result.Reasons, err.Error())
		return result
	}
	if digest, ok := verifyPayloadDigest(checkpoint); ok {
		if checkpoint.PayloadDigest == digest {
			result.PayloadDigestState = StatePass
		} else {
			result.PayloadDigestState = StateFail
			result.Reasons = append(result.Reasons, "checkpoint payload_digest does not match canonical payload")
		}
	} else {
		result.PayloadDigestState = StateCannotVerify
		result.Reasons = append(result.Reasons, "checkpoint payload cannot be canonicalized")
	}
	if verifySignature(checkpoint) {
		result.SignatureState = StatePass
	} else {
		result.SignatureState = StateFail
		result.Reasons = append(result.Reasons, "checkpoint signature is invalid")
	}
	compareBindings(&result, expected, checkpoint.Payload)
	applyPolicy(&result, checkpoint, policy)
	finalize(&result)
	return result
}

func VerifySet(runDir string, checkpoints []SignedCheckpoint, policy *TrustedCheckpointPolicy) VerificationResult {
	result := VerificationResult{
		SchemaVersion:        VerificationSchemaVersion,
		Result:               StatePass,
		TrustScope:           TrustScopeLocalSigned,
		SequenceState:        StatePass,
		ReplayFreshnessState: StateNotAssessed,
		SignerAuthorityState: StateNotAssessed,
		Reasons:              []string{},
	}
	if len(checkpoints) == 0 {
		result.Result = StateCannotVerify
		result.SequenceState = StateCannotVerify
		result.Reasons = append(result.Reasons, "no checkpoints supplied")
		return result
	}
	runID := checkpoints[0].RunID
	previousDigest := ""
	for i, cp := range checkpoints {
		checkpointResult := Verify(runDir, cp, policy)
		mergeSetVerification(&result, checkpointResult)
		if checkpointResult.Result == StateFail {
			result.Result = StateFail
			result.SequenceState = worseState(result.SequenceState, StateFail)
			result.Reasons = append(result.Reasons, fmt.Sprintf("checkpoint %s failed verification", cp.CheckpointID))
			return result
		}
		if checkpointResult.Result == StateCannotVerify {
			result.Result = StateCannotVerify
			result.SequenceState = worseState(result.SequenceState, StateCannotVerify)
			result.Reasons = append(result.Reasons, fmt.Sprintf("checkpoint %s cannot verify", cp.CheckpointID))
			return result
		}
		if cp.RunID != runID {
			result.Result = StateFail
			result.RunBindingState = StateFail
			result.SequenceState = StateFail
			result.Reasons = append(result.Reasons, fmt.Sprintf("checkpoint %s belongs to run %s, expected %s", cp.CheckpointID, cp.RunID, runID))
			return result
		}
		if cp.Sequence != i {
			result.Result = StateFail
			result.SequenceState = StateFail
			result.Reasons = append(result.Reasons, fmt.Sprintf("checkpoint sequence expected %d got %d", i, cp.Sequence))
			return result
		}
		if cp.Payload.PreviousCheckpointDigest != previousDigest {
			result.Result = StateFail
			result.SequenceState = StateFail
			result.Reasons = append(result.Reasons, fmt.Sprintf("checkpoint %s previous digest does not match prior checkpoint", cp.CheckpointID))
			return result
		}
		previousDigest = cp.PayloadDigest
	}
	return result
}

func baseResult(checkpoint SignedCheckpoint) VerificationResult {
	return VerificationResult{
		SchemaVersion:        VerificationSchemaVersion,
		CheckpointID:         checkpoint.CheckpointID,
		RunID:                checkpoint.RunID,
		Result:               StatePass,
		TrustScope:           TrustScopeLocalSigned,
		PayloadDigestState:   StateNotAssessed,
		SignatureState:       StateNotAssessed,
		RunBindingState:      StatePass,
		ChainBindingState:    StatePass,
		SourceBindingState:   StatePass,
		NonceBindingState:    StatePass,
		SequenceState:        StatePass,
		SignerAuthorityState: StateNotAssessed,
		ReplayFreshnessState: StateNotAssessed,
		Reasons:              []string{},
	}
}

func verifyPayloadDigest(checkpoint SignedCheckpoint) (string, bool) {
	canonical, err := trace.CanonicalJSON(checkpoint.Payload)
	if err != nil {
		return "", false
	}
	return trace.SHA256Hex(string(canonical)), true
}

func verifySignature(checkpoint SignedCheckpoint) bool {
	if checkpoint.Signature.Algorithm != SignatureAlgorithmEd25519 {
		return false
	}
	publicKey, err := base64.StdEncoding.DecodeString(checkpoint.Signature.PublicKey)
	if err != nil || len(publicKey) != ed25519.PublicKeySize {
		return false
	}
	signature, err := base64.StdEncoding.DecodeString(checkpoint.Signature.Signature)
	if err != nil || len(signature) != ed25519.SignatureSize {
		return false
	}
	canonical, err := trace.CanonicalJSON(checkpoint.Payload)
	if err != nil {
		return false
	}
	return ed25519.Verify(ed25519.PublicKey(publicKey), canonical, signature)
}

func compareBindings(result *VerificationResult, expected, actual Payload) {
	if actual.RunID != expected.RunID {
		result.RunBindingState = StateFail
		result.Reasons = append(result.Reasons, fmt.Sprintf("run_id mismatch: expected %s got %s", expected.RunID, actual.RunID))
	}
	if actual.EventChainHead != expected.EventChainHead || actual.EventCount != expected.EventCount {
		result.ChainBindingState = StateFail
		result.Reasons = append(result.Reasons, "event chain binding does not match selected run")
	}
	if actual.SourceSnapshotDigest != expected.SourceSnapshotDigest || actual.SourceSnapshotState != expected.SourceSnapshotState || actual.TaskHash != expected.TaskHash || actual.ContractDigest != expected.ContractDigest {
		result.SourceBindingState = StateFail
		result.Reasons = append(result.Reasons, "source, task, or contract binding does not match selected run")
	}
	if actual.RunNonce != expected.RunNonce {
		result.NonceBindingState = StateFail
		result.Reasons = append(result.Reasons, "run nonce binding does not match selected run")
	}
}

func applyPolicy(result *VerificationResult, checkpoint SignedCheckpoint, policy *TrustedCheckpointPolicy) {
	if policy == nil {
		result.SignerAuthorityState = StateNotAssessed
		result.Reasons = append(result.Reasons, "checkpoint signer authority policy is not assessed")
		return
	}
	for _, signer := range policy.AllowedSigners {
		if signer.SignerID != checkpoint.Signer.SignerID {
			continue
		}
		if signer.PublicKey == "" {
			result.SignerAuthorityState = StateCannotVerify
			result.Reasons = append(result.Reasons, "checkpoint signer policy missing public key binding")
			return
		}
		if signer.PublicKey != checkpoint.Signature.PublicKey {
			result.SignerAuthorityState = StateFail
			result.Reasons = append(result.Reasons, "checkpoint signer public key does not match policy")
			return
		}
		if signer.Authority != checkpoint.Signer.Authority {
			result.SignerAuthorityState = StateFail
			result.Reasons = append(result.Reasons, "checkpoint signer authority does not match policy")
			return
		}
		switch signer.Authority {
		case AuthorityLocalDevelopment:
			result.SignerAuthorityState = StatePass
			result.TrustScope = TrustScopeLocalSigned
		case AuthorityCIIsolatedJob:
			result.SignerAuthorityState = StateCannotVerify
			result.TrustScope = TrustScopeLocalSigned
			result.Reasons = append(result.Reasons, "ci isolated signer authority requires CI binding context")
		case AuthorityExternalWitness:
			result.SignerAuthorityState = StateNotIntegrated
			result.Reasons = append(result.Reasons, "external witness checkpoint authority is not integrated in Block 15")
		default:
			result.SignerAuthorityState = StateCannotVerify
			result.Reasons = append(result.Reasons, "checkpoint signer authority is unknown")
		}
		return
	}
	result.SignerAuthorityState = StateFail
	result.Reasons = append(result.Reasons, "checkpoint signer is not allowed by policy")
}

func mergeSetVerification(result *VerificationResult, checkpointResult VerificationResult) {
	result.PayloadDigestState = worseState(result.PayloadDigestState, checkpointResult.PayloadDigestState)
	result.SignatureState = worseState(result.SignatureState, checkpointResult.SignatureState)
	result.RunBindingState = worseState(result.RunBindingState, checkpointResult.RunBindingState)
	result.ChainBindingState = worseState(result.ChainBindingState, checkpointResult.ChainBindingState)
	result.SourceBindingState = worseState(result.SourceBindingState, checkpointResult.SourceBindingState)
	result.NonceBindingState = worseState(result.NonceBindingState, checkpointResult.NonceBindingState)
	result.SignerAuthorityState = worseState(result.SignerAuthorityState, checkpointResult.SignerAuthorityState)
	result.ReplayFreshnessState = worseState(result.ReplayFreshnessState, checkpointResult.ReplayFreshnessState)
	result.Reasons = append(result.Reasons, checkpointResult.Reasons...)
}

func worseState(left, right string) string {
	rank := map[string]int{
		StatePass:          0,
		StateNotAssessed:   1,
		StateNotIntegrated: 2,
		StateCannotVerify:  3,
		StateFail:          4,
		"":                 -1,
	}
	if rank[right] > rank[left] {
		return right
	}
	return left
}

func finalize(result *VerificationResult) {
	states := []string{
		result.PayloadDigestState,
		result.SignatureState,
		result.RunBindingState,
		result.ChainBindingState,
		result.SourceBindingState,
		result.NonceBindingState,
		result.SequenceState,
		result.SignerAuthorityState,
	}
	result.Result = StatePass
	for _, state := range states {
		if state == StateFail {
			result.Result = StateFail
			result.TrustScope = TrustScopeUntrustedShape
			return
		}
	}
	for _, state := range states {
		if state == StateCannotVerify {
			result.Result = StateCannotVerify
			return
		}
	}
}

func validateEnvelope(checkpoint SignedCheckpoint) error {
	if checkpoint.SchemaVersion != CheckpointSchemaVersion {
		return fmt.Errorf("unsupported checkpoint schema_version %s", checkpoint.SchemaVersion)
	}
	if checkpoint.Profile != ProfileEd25519Detached {
		return fmt.Errorf("unsupported checkpoint profile %s", checkpoint.Profile)
	}
	if checkpoint.HashAlgorithm != HashAlgorithmSHA256 {
		return fmt.Errorf("unsupported checkpoint hash_algorithm %s", checkpoint.HashAlgorithm)
	}
	if checkpoint.Canonical.Algorithm != trace.CanonicalSchemaAlgo || checkpoint.Canonical.Version != trace.CanonicalAlgoVersion {
		return errors.New("unsupported checkpoint canonicalization")
	}
	return nil
}

func runNonce(events []trace.Event) string {
	for _, event := range events {
		if event.EventType == trace.EventRecorderAttached {
			if value, ok := event.EventPayload["run_nonce"].(string); ok {
				return value
			}
		}
	}
	return ""
}

func decodePrivateKey(key KeyPair) (ed25519.PrivateKey, error) {
	if key.Algorithm != SignatureAlgorithmEd25519 {
		return nil, fmt.Errorf("unsupported key algorithm %s", key.Algorithm)
	}
	decoded, err := base64.StdEncoding.DecodeString(key.PrivateKey)
	if err != nil {
		return nil, err
	}
	if len(decoded) != ed25519.PrivateKeySize {
		return nil, fmt.Errorf("ed25519 private key length = %d", len(decoded))
	}
	return ed25519.PrivateKey(decoded), nil
}

func validateSequenceLink(sequence int, previousDigest string) error {
	if sequence < 0 {
		return errors.New("checkpoint sequence must be >= 0")
	}
	if sequence == 0 && previousDigest != "" {
		return errors.New("sequence 0 checkpoint must not declare previous_checkpoint_digest")
	}
	if sequence > 0 && previousDigest == "" {
		return errors.New("sequence > 0 checkpoint requires previous_checkpoint_digest")
	}
	return nil
}

func validateKeyPair(key KeyPair, privateKey ed25519.PrivateKey) error {
	if key.PublicKey == "" {
		return nil
	}
	decoded, err := base64.StdEncoding.DecodeString(key.PublicKey)
	if err != nil {
		return err
	}
	if len(decoded) != ed25519.PublicKeySize {
		return fmt.Errorf("ed25519 public key length = %d", len(decoded))
	}
	derived := privateKey.Public().(ed25519.PublicKey)
	if string(decoded) != string(derived) {
		return errors.New("private key does not match public key")
	}
	return nil
}
