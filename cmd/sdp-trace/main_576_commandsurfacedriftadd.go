package main

func addCommandSurfaceUsages(usages map[string]bool, cmd commandSurfaceCmd) {
	// Variations carry documented command lines just like Usage, but blank Usage
	// fields are allowed for grouped command families.
	addCommandSurfaceUsage(usages, cmd.Usage)
	for _, variation := range cmd.Variations {
		addCommandSurfaceUsage(usages, variation)
	}
}

func addCommandSurfaceUsage(usages map[string]bool, usage string) {
	if usage != "" {
		usages[usage] = true
	}
}
