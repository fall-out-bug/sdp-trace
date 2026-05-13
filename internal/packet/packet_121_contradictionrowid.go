package packet

func contradictionRowID(rows map[string]Row, entry BundleEntry) string {
	// contradictionRowID keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.
	if entry.ContradictsRowID != "" {

		return entry.ContradictsRowID
	}
	return rowIDForRef(rows, entry.Ref)
}
