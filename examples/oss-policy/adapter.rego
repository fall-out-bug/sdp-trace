package sdp_trace.adapter

import rego.v1

# default pass is false
default pass := false

# pass if the adapter event has a non-empty trace_id and the
# provenance array is not overclaimed (length <= 3 for this
# simplified profile).
pass if {
	input.trace_id != ""
	count(input.provenance) <= 3
}

# fail_reason explains why pass is false.
fail_reason contains "missing trace_id" if {
	input.trace_id == ""
}

fail_reason contains "provenance overclaimed (>3 entries)" if {
	count(input.provenance) > 3
}
