package trace

import (
	"errors"
	"fmt"
	"strings"
)

func firstValidationError(errs ...error) error {
	// Deterministic first-error reporting keeps CLI and test output stable during
	// artifact repair.
	for _, err := range errs {
		if err != nil {
			return err
		}
	}
	return nil
}

func requiredString(value, message string) error {
	if strings.TrimSpace(value) == "" {
		// Required trace identifiers reject whitespace-only values before they
		// can become ambiguous file names or event-chain references.
		return errors.New(message)
	}
	return nil
}

func nonNegative(value int, message string) error {
	if value < 0 {
		// Negative counts and sequence-like values cannot be reconciled with an
		// append-only trace.
		return errors.New(message)
	}
	return nil
}

func validSequence(sequence int) error {
	if sequence < 0 {
		// Event sequence numbers are append-only positions and must never move
		// below the genesis boundary.
		return fmt.Errorf("invalid sequence %d", sequence)
	}
	return nil
}
