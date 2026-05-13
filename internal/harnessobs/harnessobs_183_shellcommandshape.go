package harnessobs

func shellCommandShape(command []string) bool {
	return len(command) >= 3 && command[1] == "-c"
}
