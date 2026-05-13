package main

func doctorEnvironmentChecks() []doctorCheck {
	// Environment checks describe what the local process can know; they are not
	// external witness evidence.
	return []doctorCheck{
		{
			ID:     "local_process",
			State:  "pass",
			Reason: "current process can inspect local environment",
		},
		{
			ID:     "offline_development",
			State:  "offline_dev",
			Reason: "external CI identity is not required for local preview or wrapper readiness",
		},
	}
}
