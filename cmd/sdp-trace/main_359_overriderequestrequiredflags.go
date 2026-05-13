package main

var overrideRequestRequiredFlags = []requiredCLIFlag{
	{"out", "override request requires --out"},
	{"id", "override request requires --id"},
	{"by", "override request requires --by"},
	{"reason", "override request requires --reason"},
	{"source-ref", "override request requires --source-ref"},
	{"scope", "override request requires --scope"},
}
