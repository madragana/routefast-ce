package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aovidi/routefast-linux/internal/actuation"
	"github.com/aovidi/routefast-linux/internal/explain"
	"github.com/aovidi/routefast-linux/internal/policy"
	"github.com/aovidi/routefast-linux/internal/reasoning"
	"github.com/aovidi/routefast-linux/internal/telemetry"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := policy.Load("configs/policy.yaml")
	if err != nil {
		log.Fatalf("load policy: %v", err)
	}

	events := make(chan telemetry.Event, 64)
	go func() {
		events <- telemetry.Event{Source: "suricata", Type: "attack_suspected", Timestamp: time.Now().UTC(), Attrs: map[string]interface{}{"src_ip": "198.51.100.10"}}
		time.Sleep(3 * time.Second)
		events <- telemetry.Event{Source: "probe", Type: "bearer_failover_candidate", Timestamp: time.Now().UTC(), Attrs: map[string]interface{}{"bearer": "radio0"}}
	}()

	engine := reasoning.NewEngine()
	audit := explain.NewAudit("routefast-audit.jsonl")
	exec := actuation.NewExecutor()

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)

	batch := make([]telemetry.Event, 0, 8)
	tick := time.NewTicker(2 * time.Second)
	defer tick.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-sigc:
			cancel()
		case ev := <-events:
			batch = append(batch, ev)
		case <-tick.C:
			if len(batch) == 0 {
				continue
			}
			decision, err := engine.Decide(ctx, batch)
			if err != nil {
				log.Printf("decision error: %v", err)
				batch = batch[:0]
				continue
			}
			if err := policy.Authorize(cfg, decision.Action, decision.RollbackTTLSeconds); err != nil {
				_ = audit.Write(decision, "denied", err.Error())
				log.Printf("policy denied: %v", err)
				batch = batch[:0]
				continue
			}
			if err := exec.Apply(ctx, decision); err != nil {
				_ = audit.Write(decision, "allowed", err.Error())
				log.Printf("apply failed: %v", err)
				batch = batch[:0]
				continue
			}
			_ = audit.Write(decision, "allowed", "applied")
			batch = batch[:0]
		}
	}
}
