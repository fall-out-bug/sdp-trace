package main

func sf(name, desc string) flagMeta {
	return flagMeta{Name: name, Type: "string", Description: desc}
}
