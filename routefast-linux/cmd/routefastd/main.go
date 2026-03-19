package main

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"
    "os/signal"
    "sync"
    "syscall"
    "time"

    "github.com/aovidi/routefast-linux/internal/actuation"
    "github.com/aovidi/routefast-linux/internal/config"
    "github.com/aovidi/routefast-linux/internal/explain"
    "github.com/aovidi/routefast-linux/internal/httpx"
    "github.com/aovidi/routefast-linux/internal/lip4d"
    "github.com/aovidi/routefast-linux/internal/policy"
    "github.com/aovidi/routefast-linux/internal/reasoning"
    "github.com/aovidi/routefast-linux/internal/telemetry"
)

type server struct {
    cfg      config.Policy
    engine   *reasoning.Engine
    audit    *explain.Audit
    events   []telemetry.Event
    mu       sync.Mutex
    sequence int64
}

func main() {
    cfg, err := config.LoadPolicy("configs/policy.yaml")
    if err != nil { log.Fatal(err) }
    bearers, _ := config.LoadBearers("configs/bearers.yaml")
    s := &server{cfg: cfg, engine: reasoning.NewEngine(bearers), audit: explain.NewAudit(cfg.AuditPath), sequence: 1}

    ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
    defer cancel()

    mux := http.NewServeMux()
    mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) { _, _ = w.Write([]byte("ok")) })
    mux.HandleFunc("/event", s.handleEvent)
    mux.HandleFunc("/decision", s.handleDecision)

    srv := &http.Server{Addr: cfg.ListenAddr, Handler: mux}
    go func() {
        <-ctx.Done()
        _ = srv.Shutdown(context.Background())
    }()
    log.Printf("routefastd listening on %s", cfg.ListenAddr)
    log.Fatal(srv.ListenAndServe())
}

func (s *server) handleEvent(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost { http.Error(w, "method not allowed", http.StatusMethodNotAllowed); return }
    var ev telemetry.Event
    if err := json.NewDecoder(r.Body).Decode(&ev); err != nil { http.Error(w, err.Error(), http.StatusBadRequest); return }
    if ev.Timestamp.IsZero() { ev.Timestamp = time.Now().UTC() }
    s.mu.Lock(); s.events = append(s.events, ev); s.mu.Unlock()
    w.WriteHeader(http.StatusAccepted)
}

func (s *server) handleDecision(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost && r.Method != http.MethodGet { http.Error(w, "method not allowed", http.StatusMethodNotAllowed); return }
    s.mu.Lock()
    batch := append([]telemetry.Event(nil), s.events...)
    s.events = nil
    s.mu.Unlock()

    decision, err := s.engine.Decide(r.Context(), batch)
    if err != nil { http.Error(w, err.Error(), http.StatusInternalServerError); return }
    if err := policy.Authorize(s.cfg, decision.Action, decision.RollbackTTLSeconds); err != nil {
        _ = s.audit.Write(decision, "denied", err.Error())
        http.Error(w, err.Error(), http.StatusForbidden)
        return
    }
    if decision.PeerCorroborationRequired || policy.IsHighImpact(s.cfg, decision.Action) {
        env := s.makeEnvelope(decision)
        var validation lip4d.ValidationResponse
        if err := httpx.PostJSON(r.Context(), s.cfg.LIP4DURL+"/intent", env, &validation); err != nil {
            _ = s.audit.Write(decision, "corroboration_failed", err.Error())
            http.Error(w, err.Error(), http.StatusBadGateway)
            return
        }
    }
    var applyRes actuation.Result
    if err := httpx.PostJSON(r.Context(), s.cfg.ActuatorURL+"/apply", decision, &applyRes); err != nil {
        _ = s.audit.Write(decision, "apply_failed", err.Error())
        http.Error(w, err.Error(), http.StatusBadGateway)
        return
    }
    _ = s.audit.Write(decision, "applied", applyRes.Summary)
    w.Header().Set("Content-Type", "application/json")
    _ = json.NewEncoder(w).Encode(map[string]any{"decision": decision, "apply_result": applyRes})
}

func (s *server) makeEnvelope(d reasoning.Decision) lip4d.Envelope {
    now := time.Now().UTC()
    s.sequence++
    env := lip4d.Envelope{
        LIP4DVersion: "0.1",
        MessageType:  "PROPOSE_INTENT",
        IntentID:     fmt.Sprintf("intent-%d", now.UnixNano()),
        Sender:       lip4d.Sender{NodeID: s.cfg.NodeID, Domain: s.cfg.Domain, KeyID: "ed25519:" + s.cfg.NodeID + ":1"},
        Purpose:      lip4d.Purpose{Verb: stringsToVerb(d.Action), Object: "network-action", Target: d.Target},
        Context:      map[string]any{"impact": d.ImpactEstimate, "constraints": d.Constraints},
        Reason:       lip4d.Reason{Summary: d.Reason, Evidence: d.Evidence, EvidenceDigest: lip4d.EvidenceDigest(d.Evidence)},
        Authorization: lip4d.Authorization{Role: "edge-autonomy", ROEProfile: "defensive-resilience", RequestedCapability: capabilityForAction(d.Action)},
        Time:         lip4d.Timing{CreatedAt: now, NotAfter: now.Add(time.Duration(max(d.RollbackTTLSeconds, 60))*time.Second), TTLSeconds: max(d.RollbackTTLSeconds, 60), Sequence: s.sequence},
    }
    _ = lip4d.SignEnvelope(&env, lip4d.SeedToPrivateKey(s.cfg.KeySeed))
    return env
}

func stringsToVerb(action string) string {
    switch action {
    case "BLOCK_SOURCE": return "block"
    case "SHIFT_SERVICE_CLASS_TO_BEARER": return "reroute"
    default: return "explain"
    }
}

func capabilityForAction(action string) string {
    switch action {
    case "BLOCK_SOURCE": return "network.block.temporary"
    case "SHIFT_SERVICE_CLASS_TO_BEARER": return "route.reroute.temporary"
    default: return "network.explain"
    }
}

func max(a,b int) int { if a>b { return a }; return b }
