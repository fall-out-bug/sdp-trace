package packet

import (
	"fmt"
)

func (v *bundleValidator) add(format string, args ...any) {
	v.errors = append(v.errors, fmt.Sprintf(format, args...))
}
