package packet

import (
	"crypto/sha256"

	"encoding/json"
	"fmt"
)

func PacketDigest(packet Packet) string {
	// PacketDigest keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.
	clone := packet

	raw, err := json.Marshal(clone)
	if err != nil {
		return ""
	}
	sum := sha256.Sum256(raw)
	return "sha256:" + fmt.Sprintf("%x", sum[:])
}
