package releaseproof

import (
	"encoding/hex"
	"time"
)

type ProofStateBoolean struct {
	State  string `json:"state"`
	Value  *bool  `json:"value"`
	Reason string `json:"reason,omitempty"`
}

func buildVerification(input verificationInput) Verification {
	// Build the public proof in groups matching the trust boundaries above.
	result := baseVerification(input)
	applyManifestEvidence(&result, input)
	applySourceEvidence(&result, input.verificationState)
	applyExternalTrustDefaults(&result)
	applyReleaseTrustDefaults(&result)
	return result
}

func baseVerification(input verificationInput) Verification {
	// These fields identify the verifier output before any evidence-specific
	// status is attached.
	return Verification{
		ID:                       "contract-release-verification-block-18-19-source-bound",
		SchemaVersion:            SchemaVersion,
		ArtifactRole:             "verifier_output",
		TrustScope:               TrustScope,
		ReleaseVerificationState: input.state,
		VerifiedAt:               input.verificationTime.UTC().Format(time.RFC3339),
	}
}

func applyManifestEvidence(result *Verification, input verificationInput) {
	// Manifest evidence is local and byte-bound: reference, digest, signature
	// profile, identity policy, and accountability all come from one payload.
	result.ManifestRef = input.ref
	result.ManifestDigest = hex.EncodeToString(input.digest[:])
	result.ManifestDigestStatus = StatusMatched
	result.SignatureProfile = input.manifest.SigningProfile
	result.SignatureStatus = StatusNotAssessed
	result.IdentityPolicyRef = input.manifest.TrustedIdentityPolicyRef
	result.IdentityPolicyStatus = StatusNotAssessed
	result.Accountability = input.manifest.Accountability
}
