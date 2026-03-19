package main

import (
    "crypto/ed25519"
    "encoding/base64"
    "encoding/json"
    "log"
    "net/http"
    "sync"

    "github.com/aovidi/routefast-linux/internal/config"
    "github.com/aovidi/routefast-linux/internal/lip4d"
)

func main() {
    cfg, err := config.LoadPolicy("configs/policy.yaml")
    if err != nil { log.Fatal(err) }
    peers, _ := config.LoadPeers("configs/peers.yaml")
    keyMap := map[string]ed25519.PublicKey{}
    for _, p := range peers {
        if p.PublicKey == "" { continue }
        raw, err := base64.StdEncoding.DecodeString(p.PublicKey)
        if err == nil { keyMap[p.NodeID] = ed25519.PublicKey(raw) }
    }
    ownPriv := lip4d.SeedToPrivateKey(cfg.KeySeed)
    ownPub := ownPriv.Public().(ed25519.PublicKey)

    seen := struct{ mu sync.Mutex; m map[string]int64 }{m: map[string]int64{}}

    mux := http.NewServeMux()
    mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) { _, _ = w.Write([]byte("ok")) })
    mux.HandleFunc("/pubkey", func(w http.ResponseWriter, _ *http.Request) {
        _ = json.NewEncoder(w).Encode(map[string]string{"node_id": cfg.NodeID, "public_key": lip4d.PublicKeyString(ownPub)})
    })
    mux.HandleFunc("/intent", func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost { http.Error(w, "method not allowed", http.StatusMethodNotAllowed); return }
        var env lip4d.Envelope
        if err := json.NewDecoder(r.Body).Decode(&env); err != nil { http.Error(w, err.Error(), http.StatusBadRequest); return }
        if err := lip4d.ValidateEnvelope(env); err != nil { http.Error(w, err.Error(), http.StatusBadRequest); return }
        pub, ok := keyMap[env.Sender.NodeID]
        if !ok && env.Sender.NodeID == cfg.NodeID { pub = ownPub; ok = true }
        if !ok { http.Error(w, "unknown peer", http.StatusForbidden); return }
        if err := lip4d.VerifyEnvelope(env, pub); err != nil { http.Error(w, err.Error(), http.StatusForbidden); return }
        seen.mu.Lock()
        last, exists := seen.m[env.Sender.NodeID]
        if exists && env.Time.Sequence <= last {
            seen.mu.Unlock()
            http.Error(w, "replayed or stale sequence", http.StatusConflict)
            return
        }
        seen.m[env.Sender.NodeID] = env.Time.Sequence
        seen.mu.Unlock()
        _ = json.NewEncoder(w).Encode(lip4d.ValidationResponse{Accepted: true, Detail: "accepted under local policy"})
    })
    log.Printf("lip4dd listening on :8091 node=%s", cfg.NodeID)
    log.Fatal(http.ListenAndServe(":8091", mux))
}
