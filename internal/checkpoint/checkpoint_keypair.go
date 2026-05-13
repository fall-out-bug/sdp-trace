package checkpoint

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
)

func GenerateKeyPair(signerID string) (KeyPair, error) {
	// GenerateKeyPair keeps checkpoint proof evidence explicit and replay-bound.
	// Shape, digest, signature, chain, nonce, source, and signer states stay separate.
	// This helper does not upgrade local checkpoint data into external trust.

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
