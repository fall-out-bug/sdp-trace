package main

func writablePathPassCheck(id, path, okReason string) doctorCheck {
	// The temporary probe is removed immediately; doctor reports capability,
	// not an artifact to retain.
	return doctorCheck{
		ID:        id,
		State:     "pass",
		Reason:    okReason,
		Reference: path,
	}
}
