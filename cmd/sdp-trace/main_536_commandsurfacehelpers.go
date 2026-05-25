package main

func sf(name, desc string) flagMeta {
	return flagMeta{Name: name, Type: "string", Description: desc}
}

func pos(name, desc string, required bool) positionalMeta {
	return positionalMeta{Name: name, Description: desc, Required: required}
}

func reqFlags(flags ...flagMeta) []flagMeta {
	return flags
}

func optFlags(flags ...flagMeta) []flagMeta {
	return flags
}

func subs(names ...string) []string {
	return names
}

func reqPos(name, desc string) []positionalMeta {
	return []positionalMeta{pos(name, desc, true)}
}

func vars(v ...string) []string {
	return v
}
