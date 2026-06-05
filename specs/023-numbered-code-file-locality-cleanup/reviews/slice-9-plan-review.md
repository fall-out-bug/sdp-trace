# Slice 9 Plan Review: Witness Command Locality Cleanup

Status: accepted for implementation.

## Scope

Slice 9 is bounded to `cmd/sdp-trace/witness_[0-9]*_*.go` command shards.
Witness-related `core`, `gate`, `doctor`, and `assess` files remain outside
this slice.

## Decision Gate

- Simpler/Faster: a repo-wide rename sweep would remove numbers faster but
  would mix command families and make review evidence weaker.
- Blocking Edge Cases: witness command files include CLI usage behavior,
  customer-PKI required inputs, and exit-code behavior; grouping must preserve
  those contracts and MI thresholds.
- Existing Open Source: not applicable; this is a local file-locality cleanup
  with existing package APIs.

## Review

- Spec drift: no drift; the slice targets only the listed witness shards.
- Product drift: no product behavior change intended.
- Constitution drift: no package-boundary or dependency-direction change.
- CleanArch hex: same `cmd/sdp-trace` package, no new dependency.
- CleanCode/SOLID/DRY/YAGNI: accepted if grouped by command, options,
  validation, builders, output, and customer-PKI helpers without adding new
  abstractions.
- CRAP/MI: broader groupings must be split when file MI falls below 70.
