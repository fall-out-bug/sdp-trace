package checkpoint

import (
	"crypto/ed25519"
	"encoding/base64"
	"errors"
	"fmt"
)

func validateKeyPair(key KeyPair, privateKey ed25519.PrivateKey) error {
	// validateKeyPair keeps checkpoint proof evidence explicit and replay-bound.
	// Shape, digest, signature, chain, nonce, source, and signer states stay separate.
	// This helper does not upgrade local checkpoint data into external trust.
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
	// decodePublicKey keeps checkpoint proof evidence explicit and replay-bound.
	// Shape, digest, signature, chain, nonce, source, and signer states stay separate.
	// This helper does not upgrade local checkpoint data into external trust.

	decoded, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return nil, err
	}
	if len(decoded) != ed25519.PublicKeySize {

		return nil, fmt.Errorf("ed25519 public key length = %d", len(decoded))
	}
	return decoded, nil
}
