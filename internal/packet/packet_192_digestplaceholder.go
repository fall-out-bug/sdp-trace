package packet

import (
	"crypto/sha256"

	"fmt"
)

func digestPlaceholder(value string) string {
	sum := sha256.Sum256([]byte(value))
	return "sha256:" + fmt.Sprintf("%x", sum[:])
}
