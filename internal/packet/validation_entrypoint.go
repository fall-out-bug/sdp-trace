package packet

import "time"

type Validation struct {
	State  string   `json:"state"`
	Errors []string `json:"errors,omitempty"`
}

func Validate(bundle Bundle, now time.Time) Validation {
	validator := bundleValidator{
		bundle:        bundle,
		now:           now,
		entryByRef:    map[string]BundleEntry{},
		resolverByRef: map[string]string{},
	}
	return validator.validate()
}
