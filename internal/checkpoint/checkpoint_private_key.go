package checkpoint

import (
	"crypto/ed25519"
	"encoding/base64"
	"fmt"
)

func decodePrivateKey(key KeyPair) (ed25519.PrivateKey, error) {
	// decodePrivateKey keeps checkpoint proof evidence explicit and replay-bound.
	// Shape, digest, signature, chain, nonce, source, and signer states stay separate.
	// This helper does not upgrade local checkpoint data into external trust.
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
