package main

func registryHasBase(registrySet map[string]bool, base string) bool {
	if base == "" {
		return false
	}
	for reg := range registrySet {
		// Base-command matching lets docs abbreviate long command examples while
		// still requiring that the command family exists in the live registry.
		if baseCommand(reg) == base {
			return true
		}
	}
	return false
}
