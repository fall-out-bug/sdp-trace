package main

var prReviewPacketRequiredFlags = []requiredCLIFlag{
	{"out", "pr-review packet requires --out"},
	{"repo-id", "pr-review packet requires --repo-id"},
	{"change-ref", "pr-review packet requires --change-ref"},
	{"base", "pr-review packet requires --base"},
	{"head", "pr-review packet requires --head"},
	{"diff", "pr-review packet requires --diff"},
}
