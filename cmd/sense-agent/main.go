package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/aovidi/routefast-linux/internal/telemetry"
)

func main() {
	events := []telemetry.Event{
		{Source: "suricata", Type: "attack_suspected", Timestamp: time.Now().UTC(), Attrs: map[string]interface{}{"src_ip": "198.51.100.10", "signature": "ET.DDOS"}},
		{Source: "probe", Type: "bearer_failover_candidate", Timestamp: time.Now().UTC(), Attrs: map[string]interface{}{"primary": "fiber0", "fallback": "radio0"}},
	}
	for _, ev := range events {
		buf, _ := json.Marshal(ev)
		fmt.Println(string(buf))
	}
}
