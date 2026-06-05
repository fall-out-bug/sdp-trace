package packet

func (v *bundleValidator) validateMetadata() {
	v.validateSchemaMetadata()
	v.validateBundleIdentity()
	v.validatePacketDigest()
	v.validatePacketPolicyMetadata()
}
