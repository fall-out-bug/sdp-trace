package checkpoint

import (
	"crypto/ed25519"
	"encoding/base64"

	"github.com/fall_out_bug/sdp-trace/internal/trace"
)

func verifyPayloadDigest(checkpoint SignedCheckpoint) (string, bool) {
	// verifyPayloadDigest keeps checkpoint proof evidence explicit and replay-bound.
	// Shape, digest, signature, chain, nonce, source, and signer states stay separate.
	// This helper does not upgrade local checkpoint data into external trust.

	canonical, err := trace.CanonicalJSON(checkpoint.Payload)
	if err != nil {
		return "", false
	}
	return trace.SHA256Hex(string(canonical)), true
}

func verifySignature(checkpoint SignedCheckpoint) bool {
	// verifySignature keeps checkpoint proof evidence explicit and replay-bound.
	// Shape, digest, signature, chain, nonce, source, and signer states stay separate.
	// This helper does not upgrade local checkpoint data into external trust.
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
	// decodeSignature keeps checkpoint proof evidence explicit and replay-bound.
	// Shape, digest, signature, chain, nonce, source, and signer states stay separate.
	// This helper does not upgrade local checkpoint data into external trust.

	publicKey, publicOK := decodeSizedBase64(signature.PublicKey, ed25519.PublicKeySize)
	decodedSignature, signatureOK := decodeSizedBase64(signature.Signature, ed25519.SignatureSize)
	if !publicOK || !signatureOK {
		return nil, nil, false
	}
	return publicKey, decodedSignature, true
}

func decodeSizedBase64(value string, size int) ([]byte, bool) {
	// decodeSizedBase64 keeps checkpoint proof evidence explicit and replay-bound.
	// Shape, digest, signature, chain, nonce, source, and signer states stay separate.
	// This helper does not upgrade local checkpoint data into external trust.

	decoded, err := base64.StdEncoding.DecodeString(value)
	return decoded, err == nil && len(decoded) == size
}
