# routefast-linux v0.1 scaffold

Linux-first technical scaffold for a route.fast-style Community Edition prototype.

This repo implements the core flow:

Sense -> Reason -> Gate -> Coordinate -> Act -> Explain -> Verify -> Rollback

It includes:
- `routefastd`: local orchestrator
- `lip4dd`: LIP-4D peer protocol daemon
- `sense-agent`: telemetry normalizer stub
- `actuatord`: bounded Linux actuation stub
- JSON schemas for decisions and LIP-4D envelopes
- sample YAML policies, peers, and bearer configs
- `systemd` units

## Scope

This is a **safe defensive prototype scaffold**. It does not include offensive capabilities.
It focuses on:
- signed, time-bounded intent exchange
- peer corroboration
- multi-bearer resilience (fiber/LTE/radio/sat modeled as bearers)
- explainable actions and rollback hooks

## Build

```bash
go build ./...
```

## Run locally

Terminal 1:

```bash
go run ./cmd/routefastd
```

Terminal 2:

```bash
go run ./cmd/lip4dd
```

Terminal 3:

```bash
go run ./cmd/sense-agent
```

## Suggested next steps

1. Replace telemetry stubs with eBPF, Suricata EVE, and FRR watchers.
2. Replace the rule-based reasoner with a schema-constrained local model.
3. Add mTLS or QUIC transport for peer-to-peer LIP-4D.
4. Implement real rollback verification and FRR integration.
5. Add WireGuard overlay handling for rescue paths.
