package main

func previewBoundaries() []previewBoundary {
	// Preview is explicit about which observation boundaries are local,
	// unsupported, or not integrated.
	return append([]previewBoundary(nil), previewBoundaryRows...)
}
