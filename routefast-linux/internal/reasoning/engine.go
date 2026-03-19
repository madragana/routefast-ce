package reasoning

import (
    "context"
    "sort"
    "time"

    "github.com/aovidi/routefast-linux/internal/config"
    "github.com/aovidi/routefast-linux/internal/telemetry"
)

type Engine struct {
    Bearers []config.Bearer
}

func NewEngine(bearers []config.Bearer) *Engine { return &Engine{Bearers: bearers} }

func (e *Engine) Decide(_ context.Context, batch []telemetry.Event) (Decision, error) {
    d := Decision{Action: "EXPLAIN_ONLY", Confidence: 0.35, Reason: "insufficient evidence", Evidence: []string{}, ImpactEstimate: "low", CreatedAt: time.Now().UTC()}
    var attackIP string
    var attackSignals, failSignals int
    for _, ev := range batch {
        d.Evidence = append(d.Evidence, ev.Source+":"+ev.Type)
        switch ev.Type {
        case "attack_suspected", "suricata_alert", "pps_anomaly":
            attackSignals++
            if attackIP == "" {
                if ip, ok := ev.Attrs["src_ip"].(string); ok { attackIP = ip }
            }
        case "bearer_failover_candidate", "link_down", "bgp_withdraw", "jamming_suspected":
            failSignals++
        }
    }
    if attackSignals >= 2 && attackIP != "" {
        return Decision{Action: "BLOCK_SOURCE", Target: attackIP, Confidence: 0.92, Reason: "correlated hostile-source telemetry across multiple sensors", Evidence: d.Evidence, ImpactEstimate: "medium", RollbackTTLSeconds: 300, PeerCorroborationRequired: true, CreatedAt: time.Now().UTC()}, nil
    }
    if failSignals >= 1 {
        prefer := e.bestBearers([]string{"critical-control", "rescue-path"})
        return Decision{Action: "SHIFT_SERVICE_CLASS_TO_BEARER", Target: "critical-control", Constraints: map[string]interface{}{"prefer_media": prefer}, Confidence: 0.84, Reason: "degraded primary bearer and viable alternate bearer available", Evidence: d.Evidence, ImpactEstimate: "medium", RollbackTTLSeconds: 600, CreatedAt: time.Now().UTC()}, nil
    }
    return d, nil
}

func (e *Engine) bestBearers(purposes []string) []string {
    type scored struct { name string; score float64 }
    var scores []scored
    for _, b := range e.Bearers {
        if b.AdminState != "up" || b.OperState != "up" { continue }
        score := float64(100-b.Priority) + b.JamResistance*10 + b.TrustScore*10
        scores = append(scores, scored{name: b.MediaType, score: score})
    }
    sort.Slice(scores, func(i, j int) bool { return scores[i].score > scores[j].score })
    seen := map[string]bool{}
    out := []string{}
    for _, s := range scores {
        if !seen[s.name] { out = append(out, s.name); seen[s.name] = true }
    }
    if len(out) == 0 { out = []string{"radio", "lte"} }
    return out
}
