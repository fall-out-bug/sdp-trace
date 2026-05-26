package main

func stringSliceToSet(s []string) map[string]bool {
	set := map[string]bool{}
	for _, v := range s {
		set[v] = true
	}
	return set
}
