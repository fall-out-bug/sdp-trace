package packet

func hasMiniMaxComponent(components map[string]bool) bool {
	return components["minimax"] || components["minimax-m2.5"] || components["minimax-m2"]
}
