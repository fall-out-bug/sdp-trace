package main

func missingFromDocs(registry, docs []string) []string {
	docSet := stringSliceToSet(docs)
	var missing []string
	for _, r := range registry {
		if !docSet[r] {
			missing = append(missing, r)
		}
	}
	return missing
}
