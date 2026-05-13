package trace

import (
	"fmt"
	"strings"
)

// VerifyPayloadDigest validates the payload digest for the event.
func (event Event) VerifyPayloadDigest() error {
	// Empty payload_digest is tolerated for older or partial artifacts; present
	// digests must replay exactly.
	if strings.TrimSpace(event.PayloadDigest) == "" {
		return nil
	}
	synced, err := event.syncPayloadRepresentation()
	if err != nil {
		return err
	}
	return verifySyncedPayloadDigest(synced)
}

func verifySyncedPayloadDigest(event Event) error {
	// The error includes both expected and retained values so reviewers can
	// replay the mismatch without trusting the caller.
	computed, err := CanonicalEventPayloadDigest(event.Payload)
	if err != nil {
		return err
	}
	if event.PayloadDigest != computed {
		return fmt.Errorf("payload_digest mismatch: expected %s got %s", computed, event.PayloadDigest)
	}
	return nil
}
