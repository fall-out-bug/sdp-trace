package packet

import "fmt"

func (v *bundleValidator) validate() Validation {
	v.validateMetadata()
	v.indexManifest()
	v.validateRows()
	v.validateFindingsAndGaps()
	state := StatePass
	if len(v.errors) > 0 {
		state = StateFail
	}
	return Validation{State: state, Errors: v.errors}
}

func (v *bundleValidator) add(format string, args ...any) {
	v.errors = append(v.errors, fmt.Sprintf(format, args...))
}

// Finding/gap validation runs after row validation so theater findings and
// residual gaps can reuse the already indexed packet rows and manifest refs.
func (v *bundleValidator) validateFindingsAndGaps() {
	v.validateTheaterState()
	v.validateDecisionOwners()
	for _, finding := range v.bundle.Packet.TheaterFindings {
		v.validateTheaterFinding(finding)
	}
	for _, gap := range v.bundle.Packet.ResidualGaps {
		v.validateResidualGap(gap)
	}
}
