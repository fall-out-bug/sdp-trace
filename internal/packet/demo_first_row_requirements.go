package packet

func (c *demoFirstPacketChecker) requireToolGenerated() {
	if c.bundle.Packet.AuthoringMethod != AuthoringToolGenerated {
		c.add("demo first-packet gate requires tool_generated authoring_method, got %s", c.bundle.Packet.AuthoringMethod)
	}
}
