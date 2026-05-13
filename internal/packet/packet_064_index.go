package packet

func (c *demoFirstPacketChecker) index() {
	// index keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.
	for _, row := range c.bundle.Packet.Rows {
		c.rows[row.ID] = row
	}
	for _, entry := range c.bundle.Manifest.Entries {

		c.entryByRef[entry.Ref] = entry
	}
}
