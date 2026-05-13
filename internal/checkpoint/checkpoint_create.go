package checkpoint

import (
	"crypto/ed25519"

	"github.com/fall_out_bug/sdp-trace/internal/trace"
)

func Create(runDir string, options CreateOptions) (SignedCheckpoint, error) {
	// Create keeps checkpoint proof evidence explicit and replay-bound.
	// Shape, digest, signature, chain, nonce, source, and signer states stay separate.
	// This helper does not upgrade local checkpoint data into external trust.

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
