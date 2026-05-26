package main

func pos(name, desc string, required bool) positionalMeta {
	return positionalMeta{Name: name, Description: desc, Required: required}
}
