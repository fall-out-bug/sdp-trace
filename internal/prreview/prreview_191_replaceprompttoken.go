package prreview

import "strings"

func replacePromptToken(rendered string, replacement promptReplacement) string {
	return strings.ReplaceAll(rendered, "{{"+replacement.key+"}}", replacement.value)
}
