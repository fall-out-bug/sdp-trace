package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/checkpoint"
)

func signerMatchesCheckpoint(signer checkpoint.TrustedSigner, signed checkpoint.SignedCheckpoint) bool {
	return signer.SignerID == signed.Signer.SignerID && signer.Authority == signed.Signer.Authority && signer.PublicKey == signed.Signature.PublicKey
}
