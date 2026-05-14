package main

func updateDoctorExitForLocalChecks(result string, exitCode int, checks ...doctorCheck) (string, int) {
	for _, check := range checks {
		// Only control points that the local process can inspect affect the
		// offline doctor exit code.
		result, exitCode = updateDoctorExitForCheck(result, exitCode, check)
	}
	return result, exitCode
}
