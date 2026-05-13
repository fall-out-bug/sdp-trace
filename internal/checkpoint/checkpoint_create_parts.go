package checkpoint

import (
	"encoding/base64"

	"github.com/fall_out_bug/sdp-trace/internal/trace"
)

func signedCheckpoint(options CreateOptions, payload Payload, canonical []byte, signature []byte, publicKey string) SignedCheckpoint {
	// signedCheckpoint keeps checkpoint proof evidence explicit and replay-bound.
	// Shape, digest, signature, chain, nonce, source, and signer states stay separate.
	// This helper does not upgrade local checkpoint data into external trust.

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
	// checkpointCanonicalization keeps checkpoint proof evidence explicit and replay-bound.
	// Shape, digest, signature, chain, nonce, source, and signer states stay separate.
	// This helper does not upgrade local checkpoint data into external trust.

	return trace.Canonicalization{
		Algorithm: trace.CanonicalSchemaAlgo,
		Version:   trace.CanonicalAlgoVersion,
	}
}

func checkpointSignature(signature []byte, publicKey string) Signature {
	// checkpointSignature keeps checkpoint proof evidence explicit and replay-bound.
	// Shape, digest, signature, chain, nonce, source, and signer states stay separate.
	// This helper does not upgrade local checkpoint data into external trust.

	return Signature{
		Algorithm: SignatureAlgorithmEd25519,
		Signature: base64.StdEncoding.EncodeToString(signature),
		PublicKey: publicKey,
	}
}

func checkpointSigner(signerID string) SignerIdentity {
	// checkpointSigner keeps checkpoint proof evidence explicit and replay-bound.
	// Shape, digest, signature, chain, nonce, source, and signer states stay separate.
	// This helper does not upgrade local checkpoint data into external trust.

	return SignerIdentity{
		SignerID:     signerID,
		Authority:    AuthorityLocalDevelopment,
		KeyIsolation: KeyIsolationNotAssessed,
	}
}
