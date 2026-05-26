package main

func registryHasBase(registrySet map[string]bool, base string) bool {
	if base == "" {
		return false
	}
	for reg := range registrySet {
		if baseCommand(reg) == base {
			return true
		}
	}
	return false
}
