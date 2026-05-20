package sdp_trace.adapter

import rego.v1

# default pass is false
default pass := false

# pass if the adapter event has a non-empty trace_id, provenance is an
# array, and the provenance array is not overclaimed (length <= 3 for this
# simplified profile).
pass if {
	is_string(input.trace_id)
	input.trace_id != ""
	is_array(input.provenance)
	count(input.provenance) <= 3
}

# fail_reason explains why pass is false.
fail_reason contains "missing or non-string trace_id" if {
	not is_string(input.trace_id)
}

fail_reason contains "missing trace_id" if {
	is_string(input.trace_id)
	input.trace_id == ""
}

fail_reason contains "provenance is not an array" if {
	not is_array(input.provenance)
}

fail_reason contains "provenance overclaimed (>3 entries)" if {
	is_array(input.provenance)
	count(input.provenance) > 3
}
