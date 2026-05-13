package packet

func (c *demoFirstPacketChecker) requireUsableRowEvidenceRefs(rowID string, refs []string) {
	// requireUsableRowEvidenceRefs keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.
	for _, ref := range refs {
		entry, ok := c.entryByRef[ref]
		if !ok {

			continue
		}
		if !demoUsableEntry(entry, c.now) {
			c.add("demo first-packet gate requires %s evidence ref %q to be retained and usable", rowID, ref)
		}
	}
}
