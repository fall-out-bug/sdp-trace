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

	privateKey, err := validatedPrivateKey(options)
	if err != nil {
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

	publicKey := publicKeyForCheckpoint(options.Key.PublicKey, privateKey)
	return signedCheckpoint(options, payload, canonical, signature, publicKey), nil
}

func signedCheckpoint(options CreateOptions, payload Payload, canonical []byte, signature []byte, publicKey string) SignedCheckpoint {

	return SignedCheckpoint{
		SchemaVersion: CheckpointSchemaVersion,
		CheckpointID:  options.CheckpointID,
		RunID:         payload.RunID,
		Sequence:      options.Sequence,
		Profile:       ProfileEd25519Detached,
		Canonical:     checkpointCanonicalization(),
		HashAlgorithm: HashAlgorithmSHA256,
		Payload:       payload,
		PayloadDigest: trace.SHA256Hex(string(canonical)),
		Signature:     checkpointSignature(signature, publicKey),
		Signer:        checkpointSigner(options.SignerID),
	}
}

func checkpointCanonicalization() trace.Canonicalization {

	return trace.Canonicalization{
		Algorithm: trace.CanonicalSchemaAlgo,
		Version:   trace.CanonicalAlgoVersion,
	}
}

func checkpointSignature(signature []byte, publicKey string) Signature {

	return Signature{
		Algorithm: SignatureAlgorithmEd25519,
		Signature: base64.StdEncoding.EncodeToString(signature),
		PublicKey: publicKey,
	}
}

func checkpointSigner(signerID string) SignerIdentity {

	return SignerIdentity{
		SignerID:     signerID,
		Authority:    AuthorityLocalDevelopment,
		KeyIsolation: KeyIsolationNotAssessed,
	}
}

func publicKeyForCheckpoint(configured string, privateKey ed25519.PrivateKey) string {
	if configured != "" {

		return configured
	}

	return base64.StdEncoding.EncodeToString(privateKey.Public().(ed25519.PublicKey))
}

func validatedPrivateKey(options CreateOptions) (ed25519.PrivateKey, error) {

	if err := validateCreateOptions(options); err != nil {
		return nil, err
	}
	privateKey, err := decodePrivateKey(options.Key)
	if err != nil {
		return nil, err
	}
	return privateKey, validateKeyPair(options.Key, privateKey)
}

func validateCreateOptions(options CreateOptions) error {
	if strings.TrimSpace(options.CheckpointID) == "" {

		return errors.New("checkpoint_id is required")
	}
	if strings.TrimSpace(options.SignerID) == "" {

		return errors.New("signer_id is required")
	}
	return validateSequenceLink(options.Sequence, options.PreviousCheckpointDigest)
}
func BuildPayload(runDir, previousCheckpointDigest string) (Payload, error) {

	artifact, err := validatedRunArtifact(runDir)
	if err != nil {
		return Payload{}, err
	}
	nonce := runNonce(artifact.Events)
	if nonce == "" {

		return Payload{}, errors.New("run nonce missing from recorder_attached event")
	}
	return payloadFromArtifact(artifact, nonce, previousCheckpointDigest), nil
}

func payloadFromArtifact(artifact trace.RunArtifact, nonce, previousCheckpointDigest string) Payload {
	manifest := artifact.Manifest

	return Payload{
		RunID:                    manifest.RunID,
		RunNonce:                 nonce,
		EventChainHead:           manifestChainHead(manifest),
		EventCount:               manifest.EventCount,
		SourceSnapshotDigest:     manifest.SourceSnapshot,
		SourceSnapshotState:      manifest.SourceState,
		TaskHash:                 trace.SHA256Hex(manifest.Task),
		ContractDigest:           manifest.ContractDigest,
		PreviousCheckpointDigest: previousCheckpointDigest,
		ReplayContext:            notAssessedReplayContext(),
	}
}

func notAssessedReplayContext() ReplayContext {

	return ReplayContext{
		Repository: "not_assessed",
		Ref:        "not_assessed",
		CommitSHA:  "not_assessed",
	}
}

func manifestChainHead(manifest trace.RunManifest) string {
	if manifest.EventChainHead != "" {
		return manifest.EventChainHead
	}

	return manifest.FinalChainHead
}

func validatedRunArtifact(runDir string) (trace.RunArtifact, error) {

	artifact, err := trace.OpenRunArtifact(runDir)
	if err != nil {
		return trace.RunArtifact{}, err
	}
	if err := artifact.Manifest.Validate(); err != nil {
		return trace.RunArtifact{}, err
	}
	if artifact.Manifest.EventCount != len(artifact.Events) {

		return trace.RunArtifact{}, fmt.Errorf("event_count mismatch: run.json=%d files=%d", artifact.Manifest.EventCount, len(artifact.Events))
	}
	return artifact, trace.ValidateEventChain(artifact.Events)
}

func Verify(runDir string, checkpoint SignedCheckpoint, policy *TrustedCheckpointPolicy) VerificationResult {
	result := baseResult(checkpoint)
	if !applyCheckpointShape(&result, checkpoint) {
		return result
	}

	expected, err := BuildPayload(runDir, checkpoint.Payload.PreviousCheckpointDigest)
	if err != nil {

		cannotVerify(&result, err.Error())
		return result
	}
	applyPayloadDigestState(&result, checkpoint)
	applySignatureState(&result, checkpoint)
	compareBindings(&result, expected, checkpoint.Payload)
	applyPolicy(&result, checkpoint, policy)
	finalize(&result)
	return result
}

func applyCheckpointShape(result *VerificationResult, checkpoint SignedCheckpoint) bool {
	if err := validateEnvelope(checkpoint); err != nil {

		failShape(result, err.Error())
		return false
	}
	if err := validateSequenceLink(checkpoint.Sequence, checkpoint.Payload.PreviousCheckpointDigest); err != nil {

		result.SequenceState = StateFail
		failShape(result, err.Error())
		return false
	}
	return true
}

func failShape(result *VerificationResult, reason string) {

	result.Result = StateFail
	result.TrustScope = TrustScopeUntrustedShape
	result.Reasons = append(result.Reasons, reason)
}

func cannotVerify(result *VerificationResult, reason string) {

	result.Result = StateCannotVerify
	result.Reasons = append(result.Reasons, reason)
}

func applyPayloadDigestState(result *VerificationResult, checkpoint SignedCheckpoint) {
	digest, ok := verifyPayloadDigest(checkpoint)
	if !ok {

		result.PayloadDigestState = StateCannotVerify
		result.Reasons = append(result.Reasons, "checkpoint payload cannot be canonicalized")
		return
	}
	if checkpoint.PayloadDigest == digest {

		result.PayloadDigestState = StatePass
		return
	}
	result.PayloadDigestState = StateFail
	result.Reasons = append(result.Reasons, "checkpoint payload_digest does not match canonical payload")
}

func applySignatureState(result *VerificationResult, checkpoint SignedCheckpoint) {
	if verifySignature(checkpoint) {

		result.SignatureState = StatePass
		return
	}
	result.SignatureState = StateFail
	result.Reasons = append(result.Reasons, "checkpoint signature is invalid")
}

func VerifySet(runDir string, checkpoints []SignedCheckpoint, policy *TrustedCheckpointPolicy) VerificationResult {

	result := baseSetResult()
	if len(checkpoints) == 0 {
		return emptySetResult(result)
	}
	runID := checkpoints[0].RunID
	previousDigest := ""
	for i, cp := range checkpoints {

		if stop := verifySetCheckpoint(&result, runDir, cp, policy, setLinkExpectation{runID: runID, sequence: i, previousDigest: previousDigest}); stop {
			return result
		}
		previousDigest = cp.PayloadDigest
	}
	return result
}

func baseSetResult() VerificationResult {

	return VerificationResult{
		SchemaVersion:        VerificationSchemaVersion,
		Result:               StatePass,
		TrustScope:           TrustScopeLocalSigned,
		SequenceState:        StatePass,
		ReplayFreshnessState: StateNotAssessed,
		SignerAuthorityState: StateNotAssessed,
		Reasons:              []string{},
	}
}

func emptySetResult(result VerificationResult) VerificationResult {

	result.Result = StateCannotVerify
	result.SequenceState = StateCannotVerify
	result.Reasons = append(result.Reasons, "no checkpoints supplied")
	return result
}

type setLinkExpectation struct {
	runID          string
	sequence       int
	previousDigest string
}

func verifySetCheckpoint(result *VerificationResult, runDir string, cp SignedCheckpoint, policy *TrustedCheckpointPolicy, expected setLinkExpectation) bool {

	checkpointResult := Verify(runDir, cp, policy)
	mergeSetVerification(result, checkpointResult)
	if checkpointResult.Result == StateFail || checkpointResult.Result == StateCannotVerify {
		applySetCheckpointResult(result, cp, checkpointResult.Result)
		return true
	}
	return applySetLinkChecks(result, cp, expected)
}

func applySetCheckpointResult(result *VerificationResult, cp SignedCheckpoint, state string) {

	result.Result = state
	result.SequenceState = worseState(result.SequenceState, state)
	if state == StateFail {
		result.Reasons = append(result.Reasons, fmt.Sprintf("checkpoint %s failed verification", cp.CheckpointID))
		return
	}
	result.Reasons = append(result.Reasons, fmt.Sprintf("checkpoint %s cannot verify", cp.CheckpointID))
}

func applySetLinkChecks(result *VerificationResult, cp SignedCheckpoint, expected setLinkExpectation) bool {

	if setRunMismatch(result, cp, expected.runID) {
		return true
	}
	if setSequenceMismatch(result, cp, expected.sequence) {
		return true
	}
	return setPreviousDigestMismatch(result, cp, expected.previousDigest)
}

func setRunMismatch(result *VerificationResult, cp SignedCheckpoint, runID string) bool {
	if cp.RunID == runID {
		return false
	}

	result.Result = StateFail
	result.RunBindingState = StateFail
	result.SequenceState = StateFail
	result.Reasons = append(result.Reasons, fmt.Sprintf("checkpoint %s belongs to run %s, expected %s", cp.CheckpointID, cp.RunID, runID))
	return true
}

func setSequenceMismatch(result *VerificationResult, cp SignedCheckpoint, sequence int) bool {
	if cp.Sequence == sequence {
		return false
	}

	result.Result = StateFail
	result.SequenceState = StateFail
	result.Reasons = append(result.Reasons, fmt.Sprintf("checkpoint sequence expected %d got %d", sequence, cp.Sequence))
	return true
}

func setPreviousDigestMismatch(result *VerificationResult, cp SignedCheckpoint, previousDigest string) bool {
	if cp.Payload.PreviousCheckpointDigest != previousDigest {

		result.Result = StateFail
		result.SequenceState = StateFail
		result.Reasons = append(result.Reasons, fmt.Sprintf("checkpoint %s previous digest does not match prior checkpoint", cp.CheckpointID))
		return true
	}
	return false
}

func baseResult(checkpoint SignedCheckpoint) VerificationResult {

	result := VerificationResult{
		SchemaVersion: VerificationSchemaVersion,
		CheckpointID:  checkpoint.CheckpointID,
		RunID:         checkpoint.RunID,
		Result:        StatePass,
		TrustScope:    TrustScopeLocalSigned,
		Reasons:       []string{},
	}
	applyBaseEvidenceStates(&result)
	return result
}

func applyBaseEvidenceStates(result *VerificationResult) {

	result.PayloadDigestState = StateNotAssessed
	result.SignatureState = StateNotAssessed
	result.RunBindingState = StatePass
	result.ChainBindingState = StatePass
	result.SourceBindingState = StatePass
	result.NonceBindingState = StatePass
	result.SequenceState = StatePass
	result.SignerAuthorityState = StateNotAssessed
	result.ReplayFreshnessState = StateNotAssessed
}

func verificationStates(result *VerificationResult) []string {

	return []string{
		result.PayloadDigestState,
		result.SignatureState,
		result.RunBindingState,
		result.ChainBindingState,
		result.SourceBindingState,
		result.NonceBindingState,
		result.SequenceState,
		result.SignerAuthorityState,
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
	publicKey, signature, ok := decodeSignature(checkpoint.Signature)
	if !ok {
		return false
	}
	canonical, err := trace.CanonicalJSON(checkpoint.Payload)
	if err != nil {

		return false
	}
	return ed25519.Verify(ed25519.PublicKey(publicKey), canonical, signature)
}

func decodeSignature(signature Signature) ([]byte, []byte, bool) {

	publicKey, publicOK := decodeSizedBase64(signature.PublicKey, ed25519.PublicKeySize)
	decodedSignature, signatureOK := decodeSizedBase64(signature.Signature, ed25519.SignatureSize)
	if !publicOK || !signatureOK {
		return nil, nil, false
	}
	return publicKey, decodedSignature, true
}

func decodeSizedBase64(value string, size int) ([]byte, bool) {

	decoded, err := base64.StdEncoding.DecodeString(value)
	return decoded, err == nil && len(decoded) == size
}

func compareBindings(result *VerificationResult, expected, actual Payload) {

	compareRunBinding(result, expected, actual)
	compareChainBinding(result, expected, actual)
	compareSourceBinding(result, expected, actual)
	compareNonceBinding(result, expected, actual)
}

func compareRunBinding(result *VerificationResult, expected, actual Payload) {
	if actual.RunID == expected.RunID {
		return
	}

	result.RunBindingState = StateFail
	result.Reasons = append(result.Reasons, fmt.Sprintf("run_id mismatch: expected %s got %s", expected.RunID, actual.RunID))
}

func compareChainBinding(result *VerificationResult, expected, actual Payload) {
	if actual.EventChainHead == expected.EventChainHead && actual.EventCount == expected.EventCount {
		return
	}

	result.ChainBindingState = StateFail
	result.Reasons = append(result.Reasons, "event chain binding does not match selected run")
}

func compareSourceBinding(result *VerificationResult, expected, actual Payload) {
	if sourceBindingMatches(expected, actual) {
		return
	}

	result.SourceBindingState = StateFail
	result.Reasons = append(result.Reasons, "source, task, or contract binding does not match selected run")
}

func sourceBindingMatches(expected, actual Payload) bool {

	return actual.SourceSnapshotDigest == expected.SourceSnapshotDigest &&
		actual.SourceSnapshotState == expected.SourceSnapshotState &&
		actual.TaskHash == expected.TaskHash &&
		actual.ContractDigest == expected.ContractDigest
}

func compareNonceBinding(result *VerificationResult, expected, actual Payload) {
	if actual.RunNonce == expected.RunNonce {
		return
	}

	result.NonceBindingState = StateFail
	result.Reasons = append(result.Reasons, "run nonce binding does not match selected run")
}

func applyPolicy(result *VerificationResult, checkpoint SignedCheckpoint, policy *TrustedCheckpointPolicy) {

	applyPolicySignedAuthority(result, checkpoint, policy)
}

func applyPolicySignedAuthority(result *VerificationResult, checkpoint SignedCheckpoint, policy *TrustedCheckpointPolicy) {
	if policy == nil {

		applyMissingSignerPolicy(result)
		return
	}

	signer, found := findAllowedSigner(policy.AllowedSigners, checkpoint.Signer.SignerID)
	if !found {

		applySignerDenied(result)
		return
	}
	if !applySignerBindingPolicy(result, signer, checkpoint) {

		return
	}
	applySignerAuthorityPolicy(result, signer.Authority)
}

func applyMissingSignerPolicy(result *VerificationResult) {

	result.SignerAuthorityState = StateNotAssessed
	result.Reasons = append(result.Reasons, "checkpoint signer authority policy is not assessed")
}

func applySignerDenied(result *VerificationResult) {

	result.SignerAuthorityState = StateFail
	result.Reasons = append(result.Reasons, "checkpoint signer is not allowed by policy")
}

func applySignerBindingPolicy(result *VerificationResult, signer TrustedSigner, checkpoint SignedCheckpoint) bool {
	if signer.PublicKey == "" {

		result.SignerAuthorityState = StateCannotVerify
		result.Reasons = append(result.Reasons, "checkpoint signer policy missing public key binding")
		return false
	}
	if signer.PublicKey != checkpoint.Signature.PublicKey {

		result.SignerAuthorityState = StateFail
		result.Reasons = append(result.Reasons, "checkpoint signer public key does not match policy")
		return false
	}
	if signer.Authority != checkpoint.Signer.Authority {

		result.SignerAuthorityState = StateFail
		result.Reasons = append(result.Reasons, "checkpoint signer authority does not match policy")
		return false
	}
	return true
}

func findAllowedSigner(signers []TrustedSigner, signerID string) (TrustedSigner, bool) {
	for _, signer := range signers {
		if signer.SignerID == signerID {

			return signer, true
		}
	}
	return TrustedSigner{}, false
}

func applySignerAuthorityPolicy(result *VerificationResult, authority string) {
	state := signerAuthorityState[authority]
	if state == "" {

		result.SignerAuthorityState = StateCannotVerify
		result.Reasons = append(result.Reasons, "checkpoint signer authority is unknown")
		return
	}
	result.SignerAuthorityState = state
	if reason := signerAuthorityReason[authority]; reason != "" {

		result.Reasons = append(result.Reasons, reason)
	}
	if scope := signerAuthorityTrustScope[authority]; scope != "" {
		result.TrustScope = scope
	}
}

var signerAuthorityState = map[string]string{
	AuthorityLocalDevelopment: StatePass,
	AuthorityCIIsolatedJob:    StateCannotVerify,
	AuthorityExternalWitness:  StateNotIntegrated,
}

var signerAuthorityReason = map[string]string{
	AuthorityCIIsolatedJob:   "ci isolated signer authority requires CI binding context",
	AuthorityExternalWitness: "external witness checkpoint authority is not integrated in Block 15",
}

var signerAuthorityTrustScope = map[string]string{
	AuthorityLocalDevelopment: TrustScopeLocalSigned,
	AuthorityCIIsolatedJob:    TrustScopeLocalSigned,
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

	result.Result = StatePass
	if hasVerificationState(verificationStates(result), StateFail) {

		result.Result = StateFail
		result.TrustScope = TrustScopeUntrustedShape
		return
	}
	if hasVerificationState(verificationStates(result), StateCannotVerify) {

		result.Result = StateCannotVerify
	}
}

func hasVerificationState(states []string, target string) bool {
	for _, state := range states {
		if state == target {

			return true
		}
	}
	return false
}

func validateEnvelope(checkpoint SignedCheckpoint) error {

	return firstError(
		validateEnvelopeField(checkpoint.SchemaVersion == CheckpointSchemaVersion, "unsupported checkpoint schema_version %s", checkpoint.SchemaVersion),
		validateEnvelopeField(checkpoint.Profile == ProfileEd25519Detached, "unsupported checkpoint profile %s", checkpoint.Profile),
		validateEnvelopeField(checkpoint.HashAlgorithm == HashAlgorithmSHA256, "unsupported checkpoint hash_algorithm %s", checkpoint.HashAlgorithm),
		validateCanonicalization(checkpoint.Canonical),
	)
}

func validateEnvelopeField(ok bool, format, value string) error {
	if ok {
		return nil
	}

	return fmt.Errorf(format, value)
}

func validateCanonicalization(canonical trace.Canonicalization) error {
	if canonical.Algorithm == trace.CanonicalSchemaAlgo && canonical.Version == trace.CanonicalAlgoVersion {
		return nil
	}

	return errors.New("unsupported checkpoint canonicalization")
}

func firstError(errs ...error) error {
	for _, err := range errs {
		if err != nil {

			return err
		}
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

	return validatePreviousDigestForSequence(sequence, previousDigest)
}

func validatePreviousDigestForSequence(sequence int, previousDigest string) error {
	if sequence == 0 {

		return validateGenesisPreviousDigest(previousDigest)
	}
	if previousDigest == "" {

		return errors.New("sequence > 0 checkpoint requires previous_checkpoint_digest")
	}
	return nil
}

func validateGenesisPreviousDigest(previousDigest string) error {
	if previousDigest != "" {

		return errors.New("sequence 0 checkpoint must not declare previous_checkpoint_digest")
	}
	return nil
}

func validateKeyPair(key KeyPair, privateKey ed25519.PrivateKey) error {
	if key.PublicKey == "" {

		return nil
	}
	decoded, err := decodePublicKey(key.PublicKey)
	if err != nil {
		return err
	}
	derived := privateKey.Public().(ed25519.PublicKey)
	if string(decoded) != string(derived) {

		return errors.New("private key does not match public key")
	}
	return nil
}

func decodePublicKey(value string) ([]byte, error) {

	decoded, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return nil, err
	}
	if len(decoded) != ed25519.PublicKeySize {

		return nil, fmt.Errorf("ed25519 public key length = %d", len(decoded))
	}
	return decoded, nil
}
