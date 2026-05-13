package prreview

func ReadValidation(path string) (Validation, error) {
	var validation Validation
	return validation, readJSON(path, &validation)
}
