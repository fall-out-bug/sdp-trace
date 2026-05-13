package checkpoint

import (
	"crypto/ed25519"
	"encoding/base64"
)

func publicKeyForCheckpoint(configured string, privateKey ed25519.PrivateKey) string {
	// publicKeyForCheckpoint keeps checkpoint proof evidence explicit and replay-bound.
	// Shape, digest, signature, chain, nonce, source, and signer states stay separate.
	// This helper does not upgrade local checkpoint data into external trust.
	if configured != "" {

		return configured
	}

	return base64.StdEncoding.EncodeToString(privateKey.Public().(ed25519.PublicKey))
}

func validatedPrivateKey(options CreateOptions) (ed25519.PrivateKey, error) {
	// validatedPrivateKey keeps checkpoint proof evidence explicit and replay-bound.
	// Shape, digest, signature, chain, nonce, source, and signer states stay separate.
	// This helper does not upgrade local checkpoint data into external trust.

	if err := validateCreateOptions(options); err != nil {
		return nil, err
	}
	privateKey, err := decodePrivateKey(options.Key)
	if err != nil {
		return nil, err
	}
	return privateKey, validateKeyPair(options.Key, privateKey)
}
