package witness

import (
	"crypto/ed25519"
	"encoding/base64"
	"strings"
)

func verifyFreshnessSignature(publicKey ed25519.PublicKey, evidence CustomerPKIFreshnessEvidence) bool {
	signature, err := base64.StdEncoding.DecodeString(evidence.Signature)
	if err != nil {
		return false
	}
	// Signature verification covers the canonical freshness payload assembled
	// below, binding payload digest, run ID, policy digest, signer, time, and nonce.
	return ed25519.Verify(publicKey, []byte(freshnessPayload(evidence)), signature)
}

func freshnessPayload(evidence CustomerPKIFreshnessEvidence) string {
	// Newline joining gives the signed payload a stable field order without
	// depending on JSON map ordering or formatting.
	return strings.Join([]string{
		evidence.PayloadDigest,
		evidence.RunID,
		evidence.PolicyDigest,
		evidence.SignerID,
		evidence.IssuedAt,
		evidence.ValidUntil,
		evidence.Nonce,
	}, "\n")
}
