package policy

// SignerAllowed checks whether a signer/profile tuple is permitted.
func (v AuthorityPolicyValidator) SignerAllowed(signerID string, profileID string) bool {
	// A signer must match the requested profile and carry the signer scope;
	// other scopes in the same policy entry do not authorize release signing.
	for _, signer := range v.policy.AllowedSigners {
		if signerAllowedForProfile(signer, signerID, profileID) {
			return true
		}
	}
	return false
}

func signerAllowedForProfile(signer SignerAuthorityEntry, signerID, profileID string) bool {
	return signer.SignerID == signerID && signer.ScopeAllowed("signer") && signer.ProfileID == profileID
}
