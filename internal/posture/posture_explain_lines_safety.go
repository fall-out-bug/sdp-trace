package posture

func explainOutputSafetyLines(classes []string) []string {
	return formattedLines(classes, explainOutputSafetyLine)
}

func explainOutputSafetyLine(class string) string {
	return "output_safety absent=" + class
}
