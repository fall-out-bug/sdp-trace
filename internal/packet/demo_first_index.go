package packet

func (c *demoFirstPacketChecker) index() {
	for _, row := range c.bundle.Packet.Rows {
		c.rows[row.ID] = row
	}
	for _, entry := range c.bundle.Manifest.Entries {
		c.entryByRef[entry.Ref] = entry
	}
}
