package packet

func (c *demoFirstPacketChecker) requireToolGenerated() {
	// requireToolGenerated keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.
	if c.bundle.Packet.AuthoringMethod != AuthoringToolGenerated {

		c.add("demo first-packet gate requires tool_generated authoring_method, got %s", c.bundle.Packet.AuthoringMethod)
	}
}
