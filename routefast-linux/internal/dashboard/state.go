package dashboard

import (
    "sync"
    "time"
    "github.com/aovidi/routefast-linux/internal/swarm"
)
type Event struct { Time time.Time `json:"time"`; Type string `json:"type"`; Message string `json:"message"` }
type State struct {
    mu sync.RWMutex
    NodeID string `json:"node_id"`
    Role string `json:"role"`
    Peers []swarm.Peer `json:"peers"`
    LastQuorum any `json:"last_quorum"`
    RecentEvents []Event `json:"recent_events"`
    RemotePeerStatus map[string]map[string]any `json:"remote_peer_status"`
}
func NewState(nodeID, role string) *State { return &State{NodeID: nodeID, Role: role, RecentEvents: make([]Event,0,32), RemotePeerStatus: map[string]map[string]any{}} }
func (s *State) SetPeers(peers []swarm.Peer) { s.mu.Lock(); defer s.mu.Unlock(); s.Peers = peers }
func (s *State) UpdatePeerHeartbeat(nodeID string, ts time.Time) { s.mu.Lock(); defer s.mu.Unlock(); for i := range s.Peers { if s.Peers[i].NodeID == nodeID { s.Peers[i].LastHeartbeat = ts } } }
func (s *State) SetRemotePeerStatus(nodeID string, status map[string]any) { s.mu.Lock(); defer s.mu.Unlock(); s.RemotePeerStatus[nodeID] = status }
func (s *State) SetLastQuorum(q any) { s.mu.Lock(); defer s.mu.Unlock(); s.LastQuorum = q }
func (s *State) AddEvent(t, msg string) { s.mu.Lock(); defer s.mu.Unlock(); s.RecentEvents = append([]Event{{Time: time.Now().UTC(), Type: t, Message: msg}}, s.RecentEvents...); if len(s.RecentEvents) > 20 { s.RecentEvents = s.RecentEvents[:20] } }
func (s *State) Snapshot() map[string]any { s.mu.RLock(); defer s.mu.RUnlock(); return map[string]any{"node_id": s.NodeID, "role": s.Role, "peers": s.Peers, "last_quorum": s.LastQuorum, "recent_events": s.RecentEvents, "remote_peer_status": s.RemotePeerStatus} }
