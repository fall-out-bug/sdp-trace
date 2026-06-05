package harnessobs

import (
	"errors"
	"strings"
)

// validateValidateInputs checks required CLI-style inputs before resolving any
// path, preserving the current user-facing missing-option errors.
func validateValidateInputs(opts ValidateOptions) (string, string, string, error) {
	if err := requireValidateOptions(opts); err != nil {
		return "", "", "", err
	}

	return resolveValidateInputs(opts)
}

// requireValidateOptions rejects missing profile/run arguments before the
// filesystem safety layer adds path-specific error context.
func requireValidateOptions(opts ValidateOptions) error {
	if err := requireNonBlank(opts.ProfilePath, "harness validate requires --profile"); err != nil {
		return err
	}

	if err := requireNonBlank(opts.RunDir, "harness validate requires --run"); err != nil {
		return err
	}
	return nil
}

// requireNonBlank treats whitespace-only values as absent across validation and
// session option gates.
func requireNonBlank(value, message string) error {
	if strings.TrimSpace(value) == "" {
		return errors.New(message)
	}
	return nil
}
