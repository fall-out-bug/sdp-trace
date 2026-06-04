package packet

func (c *demoFirstPacketChecker) requirePassOrPartialRows(minimum int) {
	count := 0
	for _, row := range c.rows {
		// The demo gate treats partial evidence as usable progress, but not as
		// approval; the threshold only proves enough packet surface was exercised.
		if passOrPartial(row.State) {
			count++
		}
	}

	if count < minimum {
		c.add("demo first-packet gate requires at least %d pass or partial rows, got %d", minimum, count)
	}
}

func passOrPartial(state string) bool {
	return state == StatePass || state == StatePartial
}
