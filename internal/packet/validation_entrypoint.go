package packet

import "time"

func Validate(bundle Bundle, now time.Time) Validation {
	validator := bundleValidator{
		bundle:        bundle,
		now:           now,
		entryByRef:    map[string]BundleEntry{},
		resolverByRef: map[string]string{},
	}
	return validator.validate()
}
