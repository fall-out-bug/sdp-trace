package harnessobs

import (
	"errors"
	"strings"
)

func validateRawEventConfig(profile SessionProfile) error {
	// Raw event normalization is optional, but when configured the source and
	// format must travel together so later collection can replay the evidence.
	hasFormat := profile.RawEventFormat != ""
	hasSource := strings.TrimSpace(profile.RawEventSourcePath) != ""

	if unsupportedRawEventFormat(profile.RawEventFormat) {
		return errors.New("unsupported raw_event_format")
	}

	return validateRawEventPair(hasFormat, hasSource)
}

func validateRawEventPair(hasFormat, hasSource bool) error {
	if hasFormat == hasSource {
		return nil
	}
	if hasFormat {
		return errors.New("raw_event_source_path required for raw_event_format")
	}
	return errors.New("raw_event_format required for raw_event_source_path")
}

func unsupportedRawEventFormat(format string) bool {
	return format != "" && format != OpenCodeJSONLRawFormat
}
