package lip4d

import (
    "encoding/json"
    "net/http"
    "sync"
    "time"

    "github.com/aovidi/routefast-linux/internal/swarm"
)

type Server struct {
    Replay *ReplayCache
    Peers swarm.PeersConfig
    mu sync.Mutex
    lastSeen map[string]time.Time
    peerStates map[string]map[string]any
}
func NewServer(peers swarm.PeersConfig) *Server {
    return &Server{Replay: NewReplayCache(), Peers: peers, lastSeen: map[string]time.Time{}, peerStates: map[string]map[string]any{}}
}
func (s *Server) Register(mux *http.ServeMux) {
    mux.HandleFunc("/v1/heartbeat", s.handleHeartbeat)
    mux.HandleFunc("/v1/peerstate", s.handlePeerState)
    mux.HandleFunc("/v1/corroborate", s.handleCorroborate)
    mux.HandleFunc("/v1/status", s.handleStatus)
}
func decodeEnvelope(w http.ResponseWriter, r *http.Request) (Envelope, bool) {
    var env Envelope
    if err := json.NewDecoder(r.Body).Decode(&env); err != nil { http.Error(w, "invalid json", http.StatusBadRequest); return env, false }
    return env, true
}
func (s *Server) verifyPeerEnvelope(w http.ResponseWriter, env Envelope) bool {
    peer := s.Peers.ByNodeID(env.Sender.NodeID)
    if peer == nil { http.Error(w, "unknown peer", http.StatusForbidden); return false }
    if s.Replay.Seen(env.Sender.NodeID, env.Time.Sequence) { http.Error(w, "replay detected", http.StatusConflict); return false }
    if err := ValidateTimeWindow(env.Time, 15*time.Second); err != nil { http.Error(w, err.Error(), http.StatusUnauthorized); return false }
    if err := VerifyAgainstKeyString(env, peer.PubKey); err != nil { http.Error(w, err.Error(), http.StatusUnauthorized); return false }
    return true
}
func (s *Server) handleHeartbeat(w http.ResponseWriter, r *http.Request) {
    env, ok := decodeEnvelope(w, r); if !ok || !s.verifyPeerEnvelope(w, env) { return }
    s.mu.Lock(); s.lastSeen[env.Sender.NodeID] = time.Now().UTC(); s.mu.Unlock()
    _ = json.NewEncoder(w).Encode(map[string]any{"ok": true, "type": "heartbeat"})
}
func (s *Server) handlePeerState(w http.ResponseWriter, r *http.Request) {
    env, ok := decodeEnvelope(w, r); if !ok || !s.verifyPeerEnvelope(w, env) { return }
    bodyMap, _ := env.Body.(map[string]any)
    s.mu.Lock(); s.lastSeen[env.Sender.NodeID] = time.Now().UTC(); s.peerStates[env.Sender.NodeID] = bodyMap; s.mu.Unlock()
    _ = json.NewEncoder(w).Encode(map[string]any{"ok": true, "type": "peerstate"})
}
func (s *Server) handleCorroborate(w http.ResponseWriter, r *http.Request) {
    env, ok := decodeEnvelope(w, r); if !ok || !s.verifyPeerEnvelope(w, env) { return }
    _ = json.NewEncoder(w).Encode(map[string]any{"ok": true, "type": "corroborate", "intent_id": env.IntentID, "verdict": "agree", "confidence": 0.76, "evidence": []string{"link:degraded", "route:instability"}})
}
func (s *Server) handleStatus(w http.ResponseWriter, r *http.Request) {
    s.mu.Lock(); defer s.mu.Unlock()
    _ = json.NewEncoder(w).Encode(map[string]any{"status": "up", "peer_count": len(s.Peers.Peers), "last_seen": s.lastSeen, "peer_states": s.peerStates})
}
