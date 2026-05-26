package main

import "sort"

func sortedStringKeys(m map[string]bool) []string {
	list := make([]string, 0, len(m))
	for k := range m {
		list = append(list, k)
	}
	sort.Strings(list)
	return list
}
