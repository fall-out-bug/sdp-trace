package harnessobs

import (
	"errors"
)

func unsupportedSessionStreamCapture(mode string) error {
	// unsupportedSessionStreamCapture keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	switch mode {
	case ContentDigestOnly, ContentRetainedSafe:

		return errors.New("stream_capture mode not implemented")
	default:
		return errors.New("unsupported stream_capture")
	}
}
