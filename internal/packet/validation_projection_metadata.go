package packet

func (v *bundleValidator) validateProjectionMetadata() {
	projection := v.bundle.Packet.Projection
	if invalidCanonicalProjection(projection) {
		v.add("canonical projection must be %q", ProjectionCanonical)
	}
	if missingNonCanonicalArtifactRef(projection) {
		v.add("non-canonical packet projection requires artifact_ref for canonical uploaded packet")
	}
}
