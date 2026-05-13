package prreview

func copyDiffRef(inputDir, diffPath string) (SafeRef, error) {
	return copyInput(inputDir, "diff.patch", diffPath, RefKindDiff, ContentUnifiedDiff)
}
