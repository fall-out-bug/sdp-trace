package harnessobs

import "errors"

func normalizeSessionStreamCapture(profile *SessionProfile) error {
	// Empty capture mode defaults to disabled before rejecting unsupported
	// modes, preserving profile-load mutation order.
	defaultSessionStreamCapture(profile)
	if profile.StreamCapture == "disabled" {
		return nil
	}

	return unsupportedSessionStreamCapture(profile.StreamCapture)
}

func defaultSessionStreamCapture(profile *SessionProfile) {
	if profile.StreamCapture == "" {
		profile.StreamCapture = "disabled"
	}
}

func unsupportedSessionStreamCapture(mode string) error {
	switch mode {
	case ContentDigestOnly, ContentRetainedSafe:
		return errors.New("stream_capture mode not implemented")
	default:
		return errors.New("unsupported stream_capture")
	}
}
