package harnessobs

import "errors"

// Event content validation handles payload state and unavailable-field markers
// after identity and reference checks have accepted the event envelope.
func validateEventContent(event Event) error {
	if !validContentState(event.ContentState) {
		return errors.New("invalid content_state")
	}
	return nil
}

func validateUnavailableFields(fields []UnavailableField) error {
	for _, field := range fields {
		if !validUnavailableField(field) {
			return errors.New("invalid unavailable_fields")
		}
	}
	return nil
}

func validUnavailableField(field UnavailableField) bool {
	return safeIDPattern.MatchString(field.Field) &&
		field.State == StateNotAssessed &&
		safeIDPattern.MatchString(field.ReasonCode)
}
