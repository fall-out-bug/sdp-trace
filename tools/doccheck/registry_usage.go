package main

func addUsage(usages map[string]bool, usage string) {
	if usage != "" {
		usages[usage] = true
	}
}

func addUsages(usages map[string]bool, values []string) {
	for _, v := range values {
		addUsage(usages, v)
	}
}
