package main

import "github.com/fall_out_bug/sdp-trace/internal/checkpoint"

func policyAllowsSigner(policy checkpoint.TrustedCheckpointPolicy, signed checkpoint.SignedCheckpoint) bool {
	for _, signer := range policy.AllowedSigners {
		// Policy must bind signer id, authority class, and public key to the
		// checkpoint signature.
		if signerMatchesCheckpoint(signer, signed) {
			return true
		}
	}
	return false
}

func signerMatchesCheckpoint(signer checkpoint.TrustedSigner, signed checkpoint.SignedCheckpoint) bool {
	return signer.SignerID == signed.Signer.SignerID && signer.Authority == signed.Signer.Authority && signer.PublicKey == signed.Signature.PublicKey
}
