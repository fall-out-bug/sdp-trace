package packet

func (v *bundleValidator) validateMetadata() {
	// validateMetadata keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.

	v.validateSchemaMetadata()
	v.validateBundleIdentity()
	v.validatePacketDigest()
	v.validatePacketPolicyMetadata()
}
