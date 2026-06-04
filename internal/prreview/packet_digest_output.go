package prreview

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
)

// Packet digest helpers keep hashing canonical.
//
// The digest is computed from JSON with packet_digest cleared so callers can
// replay the value from checked-in packet JSON.
func digestJSON(value any) (string, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return "", err
	}
	digest := sha256.Sum256(data)
	return hex.EncodeToString(digest[:]), nil
}

func packetDigest(packet Packet) (string, error) {
	canonical := packet
	canonical.PacketDigest = ""
	return digestJSON(canonical)
}
