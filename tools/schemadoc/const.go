package main

const (
	indexPath         = "schema/index.json"
	schemaGlob        = "schema/*.schema.json"
	readmePath        = "schema/README.md"
	statusCurrent     = "current"
	statusHistorical  = "historical"
	statusNotAssessed = "not_assessed"
	examplePresent    = "present"
	exampleNotAssessed = "not_assessed"
	indexVersion      = "1"
)

var validStatuses = map[string]bool{
	statusCurrent:     true,
	statusHistorical:  true,
	statusNotAssessed: true,
}

var validExampleCoverages = map[string]bool{
	examplePresent:     true,
	exampleNotAssessed: true,
	"":                 true, // absent means not_assessed
}
