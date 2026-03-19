# routefast-linux v0.3

RouteFast Linux prototype with real telemetry watcher stubs.

This cut adds watcher-based ingestion for:
- Suricata EVE JSON
- FRR log files
- eBPF exporter JSONL feeds
- kernel link state polling via `ip -json link show`

The watchers emit normalized telemetry events which can be ingested by `routefastd`.
