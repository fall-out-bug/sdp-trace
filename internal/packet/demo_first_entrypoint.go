package packet

import "time"

func CheckDemoFirstPacket(bundle Bundle, now time.Time) Validation {
	validation := Validate(bundle, now)
	check := demoFirstPacketChecker{
		bundle:     bundle,
		now:        now,
		rows:       map[string]Row{},
		entryByRef: map[string]BundleEntry{},
		errors:     append([]string(nil), validation.Errors...),
	}
	return check.validate()
}
