package packet

import (
	"fmt"
)

func (c *demoFirstPacketChecker) add(format string, args ...any) {
	c.errors = append(c.errors, fmt.Sprintf(format, args...))
}
