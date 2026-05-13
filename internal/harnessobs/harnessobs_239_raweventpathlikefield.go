package harnessobs

func rawEventPathLikeField(path string, rawEvent bool) bool {
	return rawEvent && rawPathLikeField(path)
}
