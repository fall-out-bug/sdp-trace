package packet

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
