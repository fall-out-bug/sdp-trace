package packet

import (
	"errors"
	"strings"
	"time"
)

func renderablePacket(bundle Bundle) (Packet, error) {
	// renderablePacket keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.

	validation := Validate(bundle, time.Now().UTC())
	if validation.State != StatePass {

		return Packet{}, errors.New(strings.Join(validation.Errors, "; "))
	}
	return bundle.Packet, nil
}
