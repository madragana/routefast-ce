# routefast-linux v0.2 prototype

Linux-first defensive prototype for a route.fast-style router swarm.

This version is more than a scaffold:
- `routefastd` exposes an HTTP control API for telemetry and decisions
- `lip4dd` validates signed LIP-4D intent envelopes with Ed25519
- `actuatord` produces bounded Linux command plans and can run in dry-run or live mode
- `sense-agent` posts sample telemetry to drive the flow
- policy, peers, and bearer configs are file-backed and editable

## Components

- `routefastd` on `:8090`
- `lip4dd` on `:8091`
- `actuatord` on `:8092`
- `sense-agent` sends sample events to `routefastd`

## Build

```bash
go build ./...
```

## Demo

Start three terminals:

```bash
go run ./cmd/actuatord
```

```bash
go run ./cmd/lip4dd
```

```bash
go run ./cmd/routefastd
```

In a fourth terminal:

```bash
go run ./cmd/sense-agent
curl -X POST http://127.0.0.1:8090/decision
```

Expected result:
- attack telemetry becomes a `BLOCK_SOURCE` decision
- the decision is wrapped as a signed LIP-4D envelope and validated by `lip4dd`
- `actuatord` returns a dry-run `nft` command plan
- `routefast-audit.jsonl` records the decision and summary

## Live mode

By default the actuator is dry-run. To make it execute commands:

```bash
ROUTEFAST_DRY_RUN=false go run ./cmd/actuatord
```

Use live mode only on a controlled lab box.

## Next steps

- replace `sense-agent` with Suricata EVE, FRR, and eBPF watchers
- add real peer-to-peer corroboration across multiple nodes
- extend actuation with FRR and WireGuard rescue path logic
- add replay cache persistence and per-peer policy scopes
