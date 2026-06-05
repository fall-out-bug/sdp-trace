package harnessobs

import "fmt"

// Profile family validation applies the same supported-family set to required
// and optional event families while preserving first-error behavior.
func validateProfileEventFamilies(requiredEventFamilies []string, optionalEventFamilies []string) error {
	if err := validateFamilyList(requiredEventFamilies); err != nil {
		return err
	}
	return validateFamilyList(optionalEventFamilies)
}

func validateFamilyList(families []string) error {
	for _, family := range families {
		if !validFamily(family) {
			return fmt.Errorf("unsupported event family: %s", family)
		}
	}
	return nil
}
