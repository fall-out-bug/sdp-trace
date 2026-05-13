package packet

func hasContradictionTarget(entry BundleEntry, rowID string) bool {

	return entry.ContradictsRef != "" && rowID != ""
}
