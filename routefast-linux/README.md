# routefast-linux v0.4

Adds FRR actuation planning and a starter eBPF probe skeleton.

Highlights:
- `vtysh` command rendering for FRR policy changes
- optional live execution with `ROUTEFAST_EXECUTE=1`
- `bpf/flow_watch.bpf.c` starter probe
- expanded actuation planning for `BLOCK_SOURCE` and `SHIFT_SERVICE_CLASS_TO_BEARER`
