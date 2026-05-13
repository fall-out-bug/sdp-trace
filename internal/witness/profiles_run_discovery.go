package witness

import (
	"github.com/fall_out_bug/sdp-trace/internal/demo"
)

func runIDsFromRoot(runsRoot string) ([]string, error) {
	// Demo discovery defines the replayable run set for witness binding.
	runDirs, err := demo.DiscoverRunDirs(runsRoot)
	if err != nil {
		return nil, err
	}
	return runIDsFromDirs(runDirs)
}
