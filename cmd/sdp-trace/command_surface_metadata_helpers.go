package main

// Command-surface metadata helpers grouped from transitional numbered shards.
func sf(name, desc string) flagMeta {
	return flagMeta{Name: name, Type: "string", Description: desc}
}

func pos(name, desc string, required bool) positionalMeta {
	return positionalMeta{Name: name, Description: desc, Required: required}
}

func reqPos(name, desc string) []positionalMeta {
	return []positionalMeta{pos(name, desc, true)}
}
