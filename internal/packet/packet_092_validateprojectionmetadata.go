package packet

func (v *bundleValidator) validateProjectionMetadata() {
	// validateProjectionMetadata keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.
	projection := v.bundle.Packet.Projection
	if invalidCanonicalProjection(projection) {

		v.add("canonical projection must be %q", ProjectionCanonical)
	}
	if missingNonCanonicalArtifactRef(projection) {
		v.add("non-canonical packet projection requires artifact_ref for canonical uploaded packet")
	}
}
