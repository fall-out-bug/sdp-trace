package main

func reqPos(name, desc string) []positionalMeta {
	return []positionalMeta{pos(name, desc, true)}
}
