package main

type flagSet struct {
	name  string
	data  map[string]string
	bools map[string]bool
	args  []string
}
