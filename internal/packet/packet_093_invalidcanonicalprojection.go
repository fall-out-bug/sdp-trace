package packet

func invalidCanonicalProjection(projection Projection) bool {
	return projection.Canonical && projection.Kind != ProjectionCanonical
}
