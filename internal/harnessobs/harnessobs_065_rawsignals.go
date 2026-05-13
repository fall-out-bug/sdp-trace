package harnessobs

func rawSignals(value any) []string {
	return rawSignalsAt("", value)
}
