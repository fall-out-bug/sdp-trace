package checkpoint

import "github.com/fall_out_bug/sdp-trace/internal/trace"

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
