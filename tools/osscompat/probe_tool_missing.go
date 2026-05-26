package main

func missingProbeTool(p probe) bool {
	return p.NeedsTool != "" && !hasTool(p.NeedsTool)
}
