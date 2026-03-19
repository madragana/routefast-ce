package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/aovidi/routefast-linux/internal/lip4d"
)

func main() {
	env := lip4d.Envelope{
		LIP4DVersion: "0.1",
		MessageType:  "PROPOSE_INTENT",
		IntentID:     "demo-intent-001",
		Sender:       lip4d.Sender{NodeID: "rt-edge-17", Domain: "routefast.local", KeyID: "ed25519:rt-edge-17:1"},
		Purpose:      lip4d.Purpose{Verb: "reroute", Object: "prefix", Target: "203.0.113.0/24"},
		Context:      map[string]any{"risk_level": "high", "prefer_media": []string{"radio", "lte"}},
		Reason:       lip4d.Reason{Summary: "demo reroute request", Evidence: []string{"probe:radio_ok"}, EvidenceDigest: "sha256:demo"},
		Authorization: lip4d.Authorization{Role: "edge-autonomy", ROEProfile: "defensive-resilience", RequestedCapability: "route.reroute.temporary"},
		Time:         lip4d.Timing{CreatedAt: time.Now().UTC(), NotAfter: time.Now().UTC().Add(5 * time.Minute), TTLSeconds: 300, Sequence: 1},
		Signature:    "demo-signature",
	}
	if err := lip4d.ValidateEnvelope(env); err != nil {
		log.Fatalf("invalid envelope: %v", err)
	}
	buf, _ := json.MarshalIndent(env, "", "  ")
	fmt.Println(string(buf))
}
