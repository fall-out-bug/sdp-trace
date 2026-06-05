package harnessobs

import "fmt"

// Parsed event validation checks the source digest before semantic validation
// so replay integrity failures keep their precise line-level error.
func validateParsedEvent(profile Profile, event Event, line []byte, lineNo int) error {
	expected := digestLine(line)
	if event.SourceDigest != expected {
		return fmt.Errorf("source line %d: source_digest_mismatch:%s", lineNo, safeEvent(event.EventID))
	}
	if err := validateEvent(profile, event); err != nil {
		return fmt.Errorf("source line %d: %w", lineNo, err)
	}
	return nil
}

func validateEvent(profile Profile, event Event) error {
	if err := validateEventIdentity(profile, event); err != nil {
		return err
	}
	if err := validateEventRefs(event); err != nil {
		return err
	}
	if err := validateEventContent(event); err != nil {
		return err
	}
	return validateUnavailableFields(event.UnavailableFields)
}
