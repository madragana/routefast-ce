package main

import (
    "context"
    "log"
    "time"

    "github.com/aovidi/routefast-linux/internal/httpx"
    "github.com/aovidi/routefast-linux/internal/telemetry"
)

func main() {
    ctx := context.Background()
    url := "http://127.0.0.1:8090/event"
    sample := []telemetry.Event{
        {Source: "suricata", Type: "attack_suspected", Timestamp: time.Now().UTC(), Attrs: map[string]interface{}{"src_ip": "198.51.100.10"}},
        {Source: "ebpf", Type: "pps_anomaly", Timestamp: time.Now().UTC(), Attrs: map[string]interface{}{"src_ip": "198.51.100.10", "pps": 125000}},
        {Source: "probe", Type: "bearer_failover_candidate", Timestamp: time.Now().UTC(), Attrs: map[string]interface{}{"bearer": "radio0"}},
    }
    for _, ev := range sample {
        if err := httpx.PostJSON(ctx, url, ev, nil); err != nil { log.Fatal(err) }
        log.Printf("sent %s:%s", ev.Source, ev.Type)
        time.Sleep(700 * time.Millisecond)
    }
    log.Printf("now call: curl -X POST http://127.0.0.1:8090/decision")
}
